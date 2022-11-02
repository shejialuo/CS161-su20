# Design Docs

## User Register

### User name

For the requirement of the username, we have to make sure the three things:

+ Each user has a unique name.
+ Usernames are case-sensitive.
+ Username cannot be empty.

So username is a special identifier we can use to generate the UUID.

### Password

For the requirement of the password, we have to make sure the three things

+ Must not assume each user has a unique password.
+ The attacker may process a precomputed lookup table containing hashes of
common passwords downloaded from the internet.
+ Password can be empty.

In order to login again, we should store the password in the DataStore. However,
we should not store the plain text, because the attacker could do anything in the
DataStore. So the idea is that we need to generate salt to the end of the password
and use `userlib.Hash`. And we need store the salt and the hashed content into the
DataStore.

Although the attacker cannot get the plain text of the password, the attacker can
change the content in the DataStore. The attacker can simply register a new user
with a password. And observe the value in the DataStore and do retry-attack.

So we need to ensure the integrity. If the data has changed, we do not allow any login.

We store one key named `<username>_login` in KeyStore. And `<username>_login(UUID)` in
DataStore to store the salt, hashedPassword and signature. To calculate the UUID we need
to first hash the `<username>_login(UUID)` and get the first 16 byte.

At this time, we can implement `InitUser`.

## User Login

Well, we first consider about the simple cases:

+ When the username is not correct. We just find whether the `<user_name>_login` exists in
Keystore.
+ We should retrieve the salt and hashedPassword and signature from DataStore. We first check
the integrity. Which means that salt and hashedPassword doesn't be changed by the attacker. Next
we verify whether the password is correct.

It's easy.
