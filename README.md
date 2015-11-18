## chatable

A simple chat server

### How to use it

#### Dependencies
- postgresql >=9.3
- redis >=2.8

#### Run it
1. change `db/dbconf.yml` to have the correct credentials.
2. run `goose up` [1] to run all migrations.
3. create a configuration file in `cmd/` (see `development.json` for example).
4. run `godep restore` [2] to install all golang dependencies.
5. run `make local` to run the server.

### APIs
Authorization is using [jwt](http://jwt.io/) protocol.

#### `POST /api/register`
- params:
  - first_name: required
  - last_name: required
  - username: required
  - password: required
  - email: required
  - phone: required

- json:

```json
{
    'status': 'success',
    'data': [
        {
            'email': 'hello@example.com',
            'first_name': 'Bob',
            'last_name': 'Bobson',
            'phone_number': '0123456789',
            'token': {
                'access_key_id': 'AN5O6G5HYN4KFY7H2KM7XX4MKRI2RPNKUQELMQ54RITJBNS5RXKA',
                'created_at': '2015-09-03T02:08:08.558583633Z',
                'expires_at': '2015-10-03T02:08:08.558584401Z',
                'is_refreshable': true,
                'modified_at': '2015-09-03T02:08:08.5585851Z',
                'refresh_token': 'FAIMPUSXVQJAMRGBSSUJHIRRZC3ROLLOLPYT5VQOR5K6OJKFY62Q',
                'secret_access_key': '5GHV7INX2QZTG4RFE6LQMR34HWSF27IUUM55P65MFW4462SLC43A'
            },
            'username': 'bob'
        }
    ],
    'error': {
        'code': 200,
        'errors': null,
        'message': ''
    },
    'page': 1,
    'per_page': 10,
    'current_page': 1
}


```

#### `POST /api/auth_token`
- params:
  - username: required
  - password: required

- json:

```json
{
    'status': 'success',
    'data': [
        {
            'access_key_id': 'R7JKSI4JJSFAANS4UTENUMGFCC4WLF2UPXWRJHJGJ6HWELFGJMEA',
            'created_at': '2015-09-03T02:15:31.533436161Z',
            'expires_at': '2015-10-03T02:15:31.533437208Z',
            'is_refreshable': true,
            'modified_at': '2015-09-03T02:15:31.533438047Z',
            'refresh_token': 'N547NOKQVBVFAR7BEZKE5FSLOB4XN5VTOYH3TMIXHTQ4O6KJ5F4Q',
            'secret_access_key': '3KHIVVGM4YNHFS7FAZQOUXM3ASBZTL23QXMHPMMNQPHMHZ4HRAHA'
        }
    ],
    'error': {
        'code': 200,
        'errors': null,
        'message': ''
    },
    'page': 1,
    'per_page': 10,
    'current_page': 1
}
```

#### `DELETE /api/auth_token`
- authentication is required

#### `GET /api/inbox`
- authentication is required

- json

```json

{
    'status': 'success',
    'data': [
        {
            'author_username': 'bob',
            'created_at': '2015-09-02T15:49:02.996617Z',
            'latest_message': 'hello'
        }
    ],
    'error': {
        'code': 200,
        'errors': None,
        'message': ''
    },
    'page': 1,
    'per_page': 10,
    'current_page': 1
}
```

#### `GET /api/thread/[username]`
- authentication is required

- json

```json
{
    'status': 'success'
    'data': [
        {
            'author': 'bob',
            'created_at': '2015-09-02T15:49:02.994339Z',
            'message': 'hello',
            'message_type': 0,
            'recipient': 'alice'
        },
        {
            'author': 'alice',
            'created_at': '2015-09-02T15:40:02.994339Z',
            'message': 'Hi bob',
            'message_type': 0,
            'recipient': 'bob'
        },
    ],
    'error': {
        'code': 200,
        'errors': null,
        'message': ''
    },
    'page': 1,
    'per_page': 10,
    'current_page': 2
}
```

#### websocket
- endpoint: `/api/ws`
- authentication is required

- send a packet:

```json
{
    'author': 'bob',
    'recipient': 'alice',
    'message': 'hello'
}

```

- receive a packet

```json
{
    'author': 'bob',
    'recipient': 'alice',
    'message': 'hello',
    'message_type': 0,
    'created_at':'2015-09-03T03:13:12.817650473Z'
}

```

### Errors

```json
{
    'status': 'fail',
    'data': [],
    'error': {
        'code': 400,
        'errors': {'error': 'Some fileds are not unique'},
        'message': 'User error'
    },
    'page': 1,
    'per_page': 10,
    'current_page': 0
}
```

### Architecture

- Servers
  - `server A`
  - `server B`
  - `server C`

- Users
  - `bob`
  - `alice`

- Message queues
  - `shared queue`
  - `A's queue`
  - `B's queue`
  - `C's queue`

- Workflow
  1. `bob` connects to `server A`.
  2. `alice` connects to `server B`, `server C`.
  3. `bob` sends a message to `alice`.
  4. the message is pushed to the `shared queue`.
  5. `shared queue` pushes the message to `A's queue`, `B's queue`, `C's queue`
    and persists the message to database.
  6. `server A` pops from `A's queue` and drops the message because `alice` is
    not connected; `server B` pops from `B's queue` and pushes to `alice`;
    `server C` pops from `B's queue` and pushes to `alice`.


[1]. `go get bitbucket.org/liamstask/goose/cmd/goose`

[2]. `go get github.com/tools/godep`
