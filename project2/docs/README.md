# Design Docs

## Requirements

At now, the design does not consider about any security requirements.

### User

Each user has a unique username. So we could make username as
the identifier.

+ *Point 1*: Use username as the identifier.

For users, the service should provide basic authentication but the most
important thing is to provide multiple sessions. And all operations will
be sync with all the sessions. The API project required to finish is
`InitUser` and `GetUser`. When there is no sessions, we should use `InitUser`.
And when there is any sessions, we should use `GetUser`. The user
may use `GetUser` when there is no sessions. So for client,
there should be a way to know how many sessions of a user exist.

+ *Point 2*: Maintain a state to show whether there is a session existing.

### File

Different uses can have different files, so we should maintain a mapping from user
to his OWN file name

+ *Point 3*: Maintain a mapping from user to his OWN file name.

### Share

Share make things complicated. Well, the file should only has one copy.
How do we know the user's shared file name?

+ *Point 4*: Maintain a mapping from user to his SHARED file name.

And we should allow a user to send invitation to the other users from
his OWN or Shared file name and revoke.

### Summary

As you can see, the requirements for this project are easy to understand.
However, the problem is the security.
