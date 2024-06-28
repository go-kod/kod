package ranklist

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kod/kod"
)

func TestRankList(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, impl Component) {
		err := impl.Add(ctx, &AddRequest{
			Key: "key",
		})
		assert.EqualError(t, err, "validate failed: Key: 'AddRequest.Member' Error:Field validation for 'Member' failed on the 'required' tag\nKey: 'AddRequest.Score' Error:Field validation for 'Score' failed on the 'required' tag")

		err = impl.Add(ctx, &AddRequest{
			Key:    "key",
			Member: "member1",
			Score:  1,
		})
		assert.Nil(t, err)

		err = impl.Add(ctx, &AddRequest{
			Key:    "key",
			Member: "member3",
			Score:  3,
		})
		assert.Nil(t, err)

		err = impl.Add(ctx, &AddRequest{
			Key:    "key",
			Member: "member2",
			Score:  2,
		})
		assert.Nil(t, err)

		list, err := impl.RankList(ctx, &RankListRequest{
			Key: "key",
			Min: "-inf",
			Max: "+inf",
		})
		assert.Nil(t, err)
		assert.Equal(t, []string{"member3", "member2", "member1"}, list)

		list, err = impl.RankList(ctx, &RankListRequest{
			Key:    "key",
			Min:    "-inf",
			Max:    "+inf",
			Count:  1,
			Offset: 1,
		})
		assert.Nil(t, err)
		assert.Equal(t, []string{"member2"}, list)
	})
}

func BenchmarkRankList(t *testing.B) {
	kod.RunTest(t, func(ctx context.Context, impl Component) {
		err := impl.Add(ctx, &AddRequest{
			Key:    "key",
			Member: "member1",
			Score:  1,
		})
		assert.Nil(t, err)

		err = impl.Add(ctx, &AddRequest{
			Key:    "key",
			Member: "member3",
			Score:  3,
		})
		assert.Nil(t, err)

		err = impl.Add(ctx, &AddRequest{
			Key:    "key",
			Member: "member2",
			Score:  2,
		})
		assert.Nil(t, err)

		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			list, err := impl.RankList(ctx, &RankListRequest{
				Key: "key",
				Min: "-inf",
				Max: "+inf",
			})
			assert.Nil(t, err)
			assert.Equal(t, []string{"member3", "member2", "member1"}, list)
		}
	})
}
