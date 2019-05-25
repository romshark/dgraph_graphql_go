package datagen

import "time"

// StatisticsWriter represents a statistics writer interface
type StatisticsWriter interface {

	// RecordUserCreation records a user creation
	RecordUserCreation(dur time.Duration)

	// RecordPostCreation records a post creation
	RecordPostCreation(dur time.Duration)
}

// StatisticsReader represents a statistics reader interface
type StatisticsReader interface {
	Users() uint64
	Posts() uint64
	UsersCreation() time.Duration
	PostsCreation() time.Duration
	UserCreationAvg() time.Duration
	PostCreationAvg() time.Duration
}

// StatisticsRecorder represents a statistics recorder interface
type StatisticsRecorder interface {
	StatisticsWriter
	StatisticsReader
}

// stats represents the generation statistics
type stats struct {
	users              uint64
	posts              uint64
	usersCreationTotal time.Duration
	postsCreationTotal time.Duration
	userCreationAvg    time.Duration
	postCreationAvg    time.Duration
}

// NewStatisticsRecorder creates a new statistics recorder instance
func NewStatisticsRecorder() StatisticsRecorder {
	return &stats{}
}

// RecordUserCreation records a user creation
func (st *stats) RecordUserCreation(dur time.Duration) {
	// Update total
	st.users++
	st.usersCreationTotal += dur

	// Recalculate average
	st.userCreationAvg =
		time.Duration(st.usersCreationTotal) / time.Duration(st.users)
}

func (st *stats) RecordPostCreation(dur time.Duration) {
	// Update total
	st.posts++
	st.postsCreationTotal += dur

	// Recalculate average
	st.postCreationAvg =
		time.Duration(st.postsCreationTotal) / time.Duration(st.posts)
}

func (st *stats) Users() uint64 { return st.users }
func (st *stats) Posts() uint64 { return st.posts }

func (st *stats) UsersCreation() time.Duration { return st.usersCreationTotal }
func (st *stats) PostsCreation() time.Duration { return st.postsCreationTotal }

func (st *stats) UserCreationAvg() time.Duration { return st.userCreationAvg }
func (st *stats) PostCreationAvg() time.Duration { return st.postCreationAvg }
