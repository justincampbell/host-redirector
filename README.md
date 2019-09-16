# OAuth Redirector

## Overview

1. A user clicks a login button on your application.
1. The application sets the redirect URL to be the hostname of the application, or the HTTP referer [sic].
1. The user is redirected the OAuth application.
1. After the user authorizes the app, they are redirected to this OAuth Redirector.
1. The OAuth Redirector redirects the user to the last set URL.

## Setup and Usage

1. Deploy this application to a hostname accessible by end-users, setting the `OAUTH_REDIRECTOR_TOKEN` to a unique string.
1. Configure your application with the URL and token of the OAuth Redirector.
1. Change your application code to set the redirect URL before the user is redirected to the OAuth application.
1. Configure the callback URL of the OAuth application to be the URL of this application with a path of `/redirect`.

## API

### `POST /set`

Sets the redirect URL to be that of the URL given in the JSON body.

```json
{
  "url": "https://test123.example.com/auth/callback"
}
```

### `GET /redirect`

Performs a 302 redirect to the previously-set URL, maintaining the query params from the request.

If no URL is set, it returns 404.

## Race Conditions

There is a chance for a race condition if 2 or more users are simultaneously logging in to different hostnames. This is typically acceptable for infrequently-accessed test/review apps. If there is a future need for stronger redirect association with a user, the redirect URL could be stored with an IP address and/or user-agent.
