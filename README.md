# Wishy
Portal to manage wish list

## Introduction
This project create different aws lambda endpoints in order to manage a wish list. For example create the main wish list page with wishes list grouped by categories

## Technology
- [Go](https://go.dev/)
- [AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html)
- [AWS API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html)
- [MongoDB](https://www.mongodb.com/docs/)
- [Serverless](https://www.serverless.com/framework/docs)

## Deploy
This project deploy over an AWS account (specified by credential stored in github secret variables) by a github action supported by Serverless framework

## Wish management
The `/manage` page lets a user sign in and manage their categories and wishes. The `users` collection must contain an `uname` and a generated `pw` for each account, in addition to the existing `_id` and `name` fields. The login form still labels these fields “Username” and “Password”. For example:

```json
{
  "_id": { "$oid": "..." },
  "name": "Display name",
  "uname": "manager",
  "pw": "generated-password"
}
```

Passwords are compared against the value stored in MongoDB; the application does not create or reset accounts. Use generated passwords and serve the management page over HTTPS. The login cookie is HTTP-only, expires after 12 hours, and is signed using that user's stored password. Categories and wishes are scoped to the signed-in user's ID. A category cannot be deleted while wishes are still assigned to it.