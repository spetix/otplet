## Prerequisites

to run otplet on your system you need the following:

- a bar engine (i3blocks, waybar, i3status)
- gpg, gpgconf and gpgagent agent properly installed

## Installation


## Setup

### Create a token

To create a token you need:

- A seed from your authentication provider
- A keypair protected with a passphrase

To store your seed safely run the following command once:

```shell
# otplet create --url=<seed> --recipient <your keypair recipient> --store <the place to store seed>
```

Once done make sure that the seed file is readable only by current user `(0400)` and now you're ready to configure your favourite bar.

Depending on the bar the events can be managed in diffent ways:
- env var containing the event (i3blocks). In this case everything is managed properly.
- specific section in configuration. In this case there are 2 specific command that can be run one for lock action and one for unlock one.
- signal (not supported yet)

