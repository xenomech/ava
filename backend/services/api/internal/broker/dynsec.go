package broker

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"ava/pkg/wire"
)

const (
	controlTopic     = "$CONTROL/dynamic-security/v1"
	secretBytes      = 24
	clientPrefix     = "hub-"
	rolePrefix       = "hub-"
	controlPlaneRole = "ava-control-plane"
)

type dynsecCommand struct {
	Command  string       `json:"command"`
	Username string       `json:"username,omitempty"`
	Password string       `json:"password,omitempty"`
	ClientID string       `json:"clientid,omitempty"`
	RoleName string       `json:"rolename,omitempty"`
	TextName string       `json:"textname,omitempty"`
	Roles    []dynsecRole `json:"roles,omitempty"`
	ACLs     []dynsecACL  `json:"acls,omitempty"`
}

type dynsecRole struct {
	RoleName string `json:"rolename"`
}

type dynsecACL struct {
	ACLType  string `json:"acltype"`
	Topic    string `json:"topic"`
	Priority int    `json:"priority,omitempty"`
	Allow    bool   `json:"allow"`
}

type dynsecEnvelope struct {
	Commands []dynsecCommand `json:"commands"`
}

func (b *Broker) ProvisionHub(ctx context.Context, tenantSlug, hubID string) (username, password string, err error) {
	if b == nil {
		return "", "", ErrNotConnected
	}

	password, err = secret()
	if err != nil {
		return "", "", err
	}

	name := clientPrefix + hubID
	topics := wire.TopicsFor(tenantSlug, hubID)

	envelope := dynsecEnvelope{Commands: []dynsecCommand{
		{Command: "deleteClient", Username: name},
		{Command: "deleteRole", RoleName: rolePrefix + hubID},
		{
			Command:  "createRole",
			RoleName: rolePrefix + hubID,
			TextName: "Ava hub " + hubID,
			ACLs: []dynsecACL{
				{ACLType: "publishClientSend", Topic: topics.State, Allow: true},
				{ACLType: "publishClientSend", Topic: topics.Status, Allow: true},
				{ACLType: "subscribeLiteral", Topic: topics.Command, Allow: true},
				{ACLType: "subscribeLiteral", Topic: topics.Apply, Allow: true},
				{ACLType: "publishClientReceive", Topic: topics.Command, Allow: true},
				{ACLType: "publishClientReceive", Topic: topics.Apply, Allow: true},
			},
		},
		{
			Command:  "createClient",
			Username: name,
			Password: password,
			TextName: "Ava hub " + hubID,
			Roles:    []dynsecRole{{RoleName: rolePrefix + hubID}},
		},
	}}

	if err := b.control(ctx, envelope); err != nil {
		return "", "", err
	}

	return name, password, nil
}

func (b *Broker) EnsureControlPlane(ctx context.Context, username string) error {
	if b == nil {
		return ErrNotConnected
	}

	if username == "" {
		return nil
	}

	return b.control(ctx, dynsecEnvelope{Commands: []dynsecCommand{
		{Command: "removeClientRole", Username: username, RoleName: controlPlaneRole},
		{Command: "deleteRole", RoleName: controlPlaneRole},
		{
			Command:  "createRole",
			RoleName: controlPlaneRole,
			TextName: "Ava control plane",
			ACLs: []dynsecACL{
				{ACLType: "publishClientSend", Topic: "ava/+/+/cmd", Allow: true},
				{ACLType: "publishClientSend", Topic: "ava/+/+/apply", Allow: true},
			},
		},
		{Command: "addClientRole", Username: username, RoleName: controlPlaneRole},
	}})
}

func (b *Broker) RevokeHub(ctx context.Context, hubID string) error {
	if b == nil {
		return ErrNotConnected
	}

	return b.control(ctx, dynsecEnvelope{Commands: []dynsecCommand{
		{Command: "deleteClient", Username: clientPrefix + hubID},
		{Command: "deleteRole", RoleName: rolePrefix + hubID},
	}})
}

func (b *Broker) control(ctx context.Context, envelope dynsecEnvelope) error {
	payload, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("encode dynsec command: %w", err)
	}

	return b.client.Publish(ctx, controlTopic, payload, false)
}

func secret() (string, error) {
	raw := make([]byte, secretBytes)

	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate broker secret: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}
