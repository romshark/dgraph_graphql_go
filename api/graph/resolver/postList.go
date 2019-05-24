package resolver

import (
	"context"
	"strconv"

	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// PostList represents the resolver of the identically named type
type PostList struct {
	root *Resolver
}

// Size resolves PostList.size
func (rsv *PostList) Size(
	ctx context.Context,
) (int32, error) {
	var query struct {
		Node []struct {
			Count int32
		}
	}
	if err := rsv.root.str.Query(
		ctx,
		`{
			node(func: has(<posts>)) {
			  count: count(uid)
			}
		}`,
		&query,
	); err != nil {
		return 0, err
	}
	if len(query.Node) < 1 {
		return 0, nil
	}
	return query.Node[0].Count, nil
}

// Version resolves PostList.version
func (rsv *PostList) Version(
	ctx context.Context,
) (string, error) {
	var query struct {
		Version []struct {
			PostsVersion string `json:"posts.version"`
		}
	}
	if err := rsv.root.str.Query(
		ctx,
		`query PostListVersion {
			version(func: has(posts.version)) {
				posts.version
			}
		}`,
		&query,
	); err != nil {
		return "00000000000000000000000000000000", err
	}
	if len(query.Version) < 1 {
		return "00000000000000000000000000000000", nil
	}

	return query.Version[0].PostsVersion, nil
}

// List resolves PostList.list
func (rsv *PostList) List(
	ctx context.Context,
	params struct {
		First  int32
		Offset int32
	},
) ([]*Post, error) {
	// Fetch page
	var query struct {
		Posts []dgraph.Post `json:"posts"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query PostsPage(
			$first: int,
			$offset: int
		) {
			posts(
				func: has(Post.id),
				orderdesc: Post.id,
				first: $first,
				offset: $offset
			) {
				uid
				Post.id
				Post.creation
				Post.title
				Post.contents
				Post.author {
					uid
				}
			}
		}`,
		map[string]string{
			"$first":  strconv.FormatInt(int64(params.First), 10),
			"$offset": strconv.FormatInt(int64(params.Offset), 10),
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	if len(query.Posts) < 1 {
		return nil, nil
	}

	resolvers := make([]*Post, len(query.Posts))
	for i, post := range query.Posts {
		resolvers[i] = &Post{
			root:      rsv.root,
			uid:       post.UID,
			id:        post.ID,
			creation:  post.Creation,
			title:     post.Title,
			contents:  post.Contents,
			authorUID: post.Author[0].UID,
		}
	}

	return resolvers, nil
}
