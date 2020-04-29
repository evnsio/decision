# Decision

<img width="100" src="./docs/decision.png"  alt="decision"/>

A Slack integration for logging decisions in Git.

---

## How does it work?

- 1️⃣ Log a decision in Slack with `/decision` 

- 2️⃣ Fill out the dialog

    <img width="400" src="./docs/populated-modal.png"  alt="decision"/>

- 3️⃣ Log the decision 

    Depending on your configuration, this will either:

    | Commit the decision direct to your decisions repo  | Create a PR for review |
    |:-:|:-:|
    | <img width="400" src="./docs/commit-to-master-message.png"  alt="decision"/> | <img width="400" src="./docs/create-pr-message.png"  alt="decision"/> |


-  4️⃣ See you decision in Git

    After committing directly, or merging the PR, your decision is logged in the correct category folder in your repo.

    <img width="400" src="./docs/decision-record.png"  alt="decision"/>

---

## Build and Run

The latest build of decision is available at [evns/decision](https://hub.docker.com/repository/docker/evns/decision) on Docker Hub. 

If you want to build it, the easiest way is via docker.  From the root directory run:

```
docker build -t decision .
```

Settings are supplied via the following command line arguments:

```
Usage of ./decision:
  -branch string
    	The branch where decisions will be committed (default "master")
  -commit-as-prs
    	Commit decisions as Pull Requests (default false)
  -commit-author string
    	The author name to use for commits (required)
  -commit-email string
    	The author email to use for commits (required)
  -github-token string
    	Your GitHub access token (required)
  -slack-token string
    	Your Slack API token starting xoxb-... (required)
  -source-owner string
    	The owner / organisation of the repo where decisions will be committed (required)
  -source-repo string
    	The repo where decisions will be committed (required)
```

For example, to run locally:

```
docker run -p 8000:8000 evns/decision
    -slack-token=xoxb-123456789101-1234567891011-abcdefghijklmnopqrstuvwx
    -slack-signing-secret=abc123def456ghi789jkl101112mno13
    -github-token=abc123def456ghi789jkl101112mno131415pqr1
    -source-owner=evnsio
    -source-repo=decisions
    -commit-author=Chris Evans
    -commit-email=my@email.com 
    -commit-as-prs=true
```

---

## Setup and usage

1. Navigate to [https://api.slack.com/apps](https://api.slack.com/apps) and select 'Create New App'.  You can call your app whatever you like, for example 'Decision'.

1. On the `Basic Information` screen copy the `Signing Secret` - you'll need this to run the app.

1. Click `Slash Commands` > `Create New Command` and enter the following:

    - Command: `/decision`
    - Request URL: `https://<your-domain>/slash`
    - Short description: `Log a decision!`

1. Click `OAuth & Permissions` and navigte to `Bot Token Scopes`.  Add the following:

    - `chat:write`
    - `chat:write.public`

1. Still in `OAuth & Permissions` click `Install App to Workspace`.  You'll be taken through the OAuth flow, and then return to the main screen where your bot token can be found. It's the one that starts with `xoxb-...`.  

1. Click `Interactivity & Shortcuts`, toggle on Interactivity, and add the following:

    - Request URL: `https://<your-domain>/action`
    - Options URL: `https://<your-domain>/options`

Once this is complete and `decision` is running you should be able to log decisions though the `/decision` command in Slack.