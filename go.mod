module github.com/simarpreeetsingh/go-explore

go 1.26.6

replace github.com/simarpreeetsingh/go-explore/lib/linkedlist => ./lib/linkedlist

replace github.com/simarpreeetsingh/go-explore/lib/stack => ./lib/stack

replace github.com/simarpreeetsingh/go-explore/lib/queue => ./lib/queue

require (
	github.com/simarpreeetsingh/go-explore/lib/linkedlist v0.0.0-00010101000000-000000000000
	github.com/simarpreeetsingh/go-explore/lib/queue v0.0.0-00010101000000-000000000000
	github.com/simarpreeetsingh/go-explore/lib/stack v0.0.0-00010101000000-000000000000
)
