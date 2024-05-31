# International and Vague Backend

Provide an interface to:

1. Building Quotes
2. Taking out Policies
3. Products
4. People

## Configuration

This service requires the [decision engine](https://github.com/intvag/decision-engine) to be running and accessible via the address in the Environment Variable `BACKEND_DECISIONS_ENGINE`.

## Routes

| Method   | Path                      | Is Authed? | Input         | Output     | Description                                                                      |
|----------|---------------------------|------------|---------------|------------|----------------------------------------------------------------------------------|
| `GET`    | `/quote`                  | No         | None          | `Quote`    | Create a new quote without any quote items                                       |
| `GET`    | `/quote/:id`              | No         | None          | `Quote`    | Return a quote, with quote items, by ID                                          |
| `POST`   | `/quote/:id`              | No         | `QuoteInput`  | `Quote`    | Get a quote item for the requested device, adding to the quote                   |
| `DELETE` | `/quote/:id/item/:item`   | No         | None          | None       | Remove a quote item from the quote                                               |
| `GET`    | `/v1/policy`              | Yes        | None          | `[Policy]` | Return the policies for the logged in user                                       |
| `GET`    | `/v1/policy/:id`          | Yes        | None          | `Policy`   | Return a policy by ID                                                            |
| `POST`   | `/v1/policy`              | Yes        | `PolicyInput` | `Policy`   | Given a `QuoteID` passed via a `PolicyInput`, create a `Policy`                  |
| `GET`    | `/v1/policy/:id/callback` | No-ish     | None          | None       | Called by payment providers to validate payment of a policy, thus making it live |


### Authentication

People, users, and auth are all done via JWT tokens from an IDP, such as Cognito, Keycloak, Auth0, whatever.

We don't care about your auth provider, customer database, or CRM; it's none of our business. We just need enough to generate policy documents.
