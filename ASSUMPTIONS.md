# Decisions
- Included a user_id as a part of the wallet table to show that we eventually want to include an owner for each wallet
- All the ids, including waller and user ids, are UUIDs because usually banks/fintechs enforce a format for their account numbers/ids. I chose to use UUIDs as there are generators available online
- The APIs should ideally also ideally include ways to provide Ids from upstream instead of just creating its own and streaming it out but as a Phase 1 take home assignment i chose to go with simple implementations
- Idempotency keys should ideally have been created as well to prevent duplicate requests and allow for retry functionality but as an initial submission I chose to make the transactions simple
- Creating DB level locks or optimistic locking mechanism should are also something I chose to skip for now for simplicity

# Reviewers
- internal/wallet contains all the service and handler logic
  - this also contains all service level unit tests
- internal/router handles routing for requests received

# Time spent
- 8hrs total