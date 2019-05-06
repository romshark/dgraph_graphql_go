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
	reaction(id: Identifier!): Reaction
}

type Mutation {
	signIn(
		email: String!
		password: String!
	): Session!

	createUser(
		email: String!
		displayName: String!
		password: String!
	): User!

	createPost(
		author: Identifier!
		title: String!
		contents: String!
	): Post!

	createReaction(
		author: Identifier!
		subject: Identifier!
		emotion: Emotion!
		message: String!
	): Reaction!

	editPost(
		post: Identifier!
		editor: Identifier!
		newTitle: String
		newContents: String
	): Post!
}

type Session {
	key: String!
	user: User!
	creation: Time!
}

type User {
	id: Identifier!
	creation: Time!
	email: String!
	displayName: String!
	posts: [Post!]!
	sessions: [Session!]!

	# publishedReactions lists all reactions published by the user
	publishedReactions: [Reaction!]!
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
