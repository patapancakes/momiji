# Momiji
A free and simple embedded message box

## Adding to your website
You can easily add Momiji to your website by placing the iframe in its source.
```html
<iframe src="https://momiji.chat"></iframe>
```

## Usage
Your site must be verified before posts can be made on it. Verification checks for the presence of the Momiji iframe and is done when the first post is made. Alternatively, you can create a blank file at `.well-known/momiji` if the first method fails.

## User identities
Momiji creates a non-reversible identity string for each poster that changes between websites. This identity is the result of a hash function containing a secret key only known to the server, the website's domain, and the poster's IP address.
