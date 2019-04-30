package graph

var schema = `
schema {
	query: Query
	mutation: Mutation
}

type Query {
	users: [User!]!
	posts: [Post!]!

	user(id: Identifier!): User
	post(id: Identifier!): Post
}

type Mutation {
	createUser(
		email: String!
		displayName: String!
	): User!

	createPost(
		author: Identifier!
		title: String!
		contents: String!
	): Post!

	createReaction(
		emotion: Emotion!
		message: String!
	): Reaction!
}

type User {
	id: Identifier!
	creation: Time!
	email: String!
	displayName: String!
	posts: [Post!]!
}

type Post {
	id: Identifier!
	author: User!
	creation: Time!
	title: String!
	contents: String!
	reactions: [Reaction!]!
}

union ReactionSubject = Reaction | Post

type Reaction {
	id: Identifier!
	creation: Time!
	subject: ReactionSubject!
	author: User!
	emotion: Emotion!
	message: String!
	reactions: [Reaction!]!
}

enum Emotion {
	happy
	angry
	excited
	fearful
	thoughtful
}

scalar Identifier
scalar Time
`
