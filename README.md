# Decision

<img style="width: 100px;" src="./docs/decision.png"  alt="decision"/>

A Slack integration for logging decisions in Git.

---

## How does it work?

### 1️⃣ Log a decision in Slack with `/decision` 

### 2️⃣ Fill out the dialog

<img style="width: 400px;" src="./docs/populated-modal.png"  alt="decision"/>

### 3️⃣ Log the decision 

Depending on your configuration, this will either:

- Commit the decision direct to your decisions repo

    <img style="width: 400px;" src="./docs/commit-to-master-message.png"  alt="decision"/>

or

- Create a PR for review.

    <img style="width: 400px;" src="./docs/create-pr-message.png"  alt="decision"/>
    
    👇

    <img style="width: 400px;" src="./docs/pr.png"  alt="decision"/>

### 4️⃣ See you decision in Git

After the committing directly, or merging the PR, your decision is logged in the coorect category folder in your repo.

<img style="width: 400px;" src="./docs/decision-record.png"  alt="decision"/>

---

## Build and Run

The easiest way to build and run is via docker.  From the root directory, decision can be built with:

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

For example, to run the image locally:

```
docker run -p 8000:8000 decision
    -slack-token=xoxb-123456789101-1234567891011-abcdefghijklmnopqrstuvwx
    -github-token=abc123def456ghi789jkl101112mno131415pqr1
    -source-owner=evnsio
    -source-repo=decisions
    -commit-author=Chris Evans
    -commit-email=my@email.com 
    -commit-as-prs=true
```

---

## Setup and usage

