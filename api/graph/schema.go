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
	createSession(
		email: String!
		password: String!
	): Session!

	# authenticate signs the client into the
	# session identified by the given key
	authenticate(
		sessionKey: String!
	): Session!

	closeSession(
		key: String!
	): Boolean!

	closeAllSessions(
		user: Identifier!
	): [String!]!

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

	editUser(
		user: Identifier!
		editor: Identifier!
		newEmail: String
		newPassword: String
	): User!

	editReaction(
		reaction: Identifier!
		editor: Identifier!
		newMessage: String!
	): Reaction!
}

type Session {
	key: String!
	user: User!
	creation: Time!
}

type User {
	id: Identifier!
	creation: Time!
	displayName: String!
	posts: [Post!]!

	# The list of active sessions can only be accessed by the profile owner
	sessions: [Session!]!

	# The email address can only be accessed by the profile owner
	email: String!

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
