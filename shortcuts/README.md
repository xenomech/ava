# Ava shortcuts

Apple Shortcuts that drive Ava from Siri. Each is a single HTTP request authenticated by a
personal access token, so there is no login step and no password on your phone.

## Getting a token

In the app, or over the API with a signed-in session:

```sh
curl -X POST https://api-stage-df53.up.railway.app/api/v1/tokens \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <your session access token>' \
  -d '{"name":"Siri","scopes":["devices:read","devices:write"]}'
```

The response carries `value` — the only time the token is ever shown. Give it the narrowest
scopes that work: `devices:read` and `devices:write` are enough for lights and scenes.

## Using these files

The signed files in `signed/` import without enabling untrusted shortcuts. After importing, edit
the one action and replace:

- `ava_pat_PASTE_YOUR_TOKEN_HERE` with your token
- `PASTE_DEVICE_UUID` with a device id from `GET /api/v1/devices`

Rename the shortcut to whatever you want to say, and Siri picks it up.

## Building one by hand

One **Get Contents of URL** action is the whole shortcut:

- URL: `https://api-stage-df53.up.railway.app/api/v1/devices/<device id>/command`
- Method: **POST**
- Headers: `Authorization` = `Bearer ava_pat_...`
- Request Body: **JSON**, `trait` (Text) = `power`, `value` (**Boolean**) = true or false

`value` must be a Boolean, not the text "true" — the API's trait validation rejects a string.

For a scene, GET `/api/v1/rooms/<room id>/scenes`, take the scene's `targets`, and POST them as
`{"targets": [...]}` to `/api/v1/devices/apply`. A scene describes a room that is on; to turn a
room off, apply `power: false` targets instead.
