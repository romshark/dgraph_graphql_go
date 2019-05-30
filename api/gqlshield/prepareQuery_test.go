package gqlshield

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPrepareQuery tests prepareQuery
func TestPrepareQuery(t *testing.T) {
	t.Run("noGaps", func(t *testing.T) {
		strs := []string{
			"f",
			"foo",
			"foobar",
		}
		for _, str := range strs {
			t.Run("", func(t *testing.T) {
				out, err := prepareQuery([]byte(str))
				require.NoError(t, err)
				require.Equal(t, string([]byte(str)), string(out))
			})
		}
	})

	t.Run("leadingSpace", func(t *testing.T) {
		strs := []string{
			" ",
			"\t",
			"\n",
			"  ",
			" \t\n",
		}
		words := []string{
			"f",
			"fo",
			"fooo",
		}
		for _, str := range strs {
			for _, word := range words {
				t.Run("", func(t *testing.T) {
					out, err := prepareQuery([]byte(str + word))
					require.NoError(t, err)
					require.Equal(t, string([]byte(word)), string(out))
				})
			}
		}
	})

	t.Run("singleGap", func(t *testing.T) {
		strs := []string{
			" ",
			"\t",
			"\n",
			"  ",
			" \t\n",
		}
		words := []string{
			"f",
			"fo",
			"fooo",
		}
		for _, str := range strs {
			for _, word := range words {
				t.Run("", func(t *testing.T) {
					out, err := prepareQuery([]byte(word + str + word))
					require.NoError(t, err)
					require.Equal(t, string([]byte(word+" "+word)), string(out))
				})
			}
		}
	})

	t.Run("multipleGaps", func(t *testing.T) {
		strs := []string{
			" ",
			"\t",
			"\n",
			"  ",
			" \t\n",
		}
		for _, str := range strs {
			t.Run("", func(t *testing.T) {
				out, err := prepareQuery(
					[]byte("foo" + str + "bar" + str + "baz"),
				)
				require.NoError(t, err)
				require.Equal(t, string([]byte("foo bar baz")), string(out))
			})
		}
	})

	t.Run("trailingSpaces", func(t *testing.T) {
		strs := []string{
			" ",
			"\t",
			"\n",
			"  ",
			" \t\n",
		}
		words := []string{
			"f",
			"fo",
			"fooo",
		}
		for _, str := range strs {
			for _, word := range words {
				t.Run("", func(t *testing.T) {
					in := []byte(word + str)
					out, err := prepareQuery(in)
					require.NoError(t, err)
					require.Equal(t, string([]byte(word)), string(out))
				})
			}
		}
	})

	t.Run("complex", func(t *testing.T) {
		out, err := prepareQuery(
			[]byte(" \tfoo \t\nbar  baz  \t  \n fuz \n\t\n"),
		)
		require.NoError(t, err)
		require.Equal(t, string([]byte("foo bar baz fuz")), string(out))
	})

	t.Run("embeddedString", func(t *testing.T) {
		out, err := prepareQuery(
			[]byte(" \tfoo \" \t\nbar  baz  \t  \n fuz \n\"\t\n"),
		)
		require.NoError(t, err)
		require.Equal(
			t,
			string([]byte("foo \" \t\nbar  baz  \t  \n fuz \n\"")),
			string(out),
		)
	})

	t.Run("multipleEmbeddedStrings", func(t *testing.T) {
		out, err := prepareQuery(
			[]byte(" \tfoo \" \t\nbar \"  baz  \" \t  \n fuz \n\"\t\n"),
		)
		require.NoError(t, err)
		require.Equal(
			t,
			string([]byte("foo \" \t\nbar \" baz \" \t  \n fuz \n\"")),
			string(out),
		)
	})
}

func TestPrepareQueryErr(t *testing.T) {
	t.Run("unclosedStringContext", func(t *testing.T) {
		out, err := prepareQuery([]byte("foo\"bar"))
		require.Error(t, err)
		require.Nil(t, out)
	})
}

func BenchmarkPrepareQuery(b *testing.B) {
	query := []byte(
		" foo bar   baz  fuz muz daaaaaaaaaaaaz               luz   jazzz    ",
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prepareQuery(query)
	}
}
