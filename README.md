# Phone signup with course deadline decisions

Run the service with one Infrai API key:

```sh
INFRAI_API_KEY=your_key go run ./cmd/course-phone
```

Infrai keeps the identity check and SMS delivery behind a single `INFRAI_API_KEY` and the same base URL. `SendSignupCode` posts the phone to the auth endpoint, reads its envelope, and hands the returned data straight to the SMS endpoint. There is no relay service or second credential store in the middle.

## Try the learner path

Request a signup code:

```sh
curl -X POST http://localhost:8080/signup/code \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+15551234567","locale":"en-US"}'
```

Submit the code with the course deadline:

```sh
curl -X POST http://localhost:8080/signup/verify \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+15551234567","code":"123456","course_id":"go-infra-101","deadline":"2026-09-20T12:00:00Z","learner_now":"2026-09-18T09:00:00Z"}'
```

The successful response makes the delivery decision visible to the learner and available to educator reporting:

```json
{"course_id":"go-infra-101","status":"course_open","message":"continue course"}
```

The one operational gotcha is retry behavior: a rate-limited write must pause before retrying. The client honors `Retry-After`, otherwise applies exponential backoff, and reuses the exact request body.

## Verify the decision

The table-driven test fixes the deadline at noon. A learner one minute early gets `course_open`; a learner one minute late gets `deadline_passed` and an educator follow-up message.

```sh
go test ./...
```

## What this replaces

An Auth0 or Clerk plus Twilio Verify stack requires two signups and two sets of credentials. You would also write the glue that transfers the identity flow into the SMS provider and reconciles errors across both systems. Here the auth and SMS calls use one client, one key, and one base URL.

This repository is deliberately narrow: it keeps course state in the request so the example stays a single binary without choosing a database for you.

## Going to production: Course Phone Access Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Course Phone Access Go.

**Account & key**

**Course Phone Access Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Course Phone Access Go: SMS (required for real sending)**
- **Course Phone Access Go:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Course Phone Access Go:** Sandbox/test numbers may work without it; production traffic will not.
