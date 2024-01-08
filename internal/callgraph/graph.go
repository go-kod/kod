package callgraph

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"sort"
)

func MakeEdgeString(src, dst string) string {
	return fmt.Sprintf("⟦%s:KoDeDgE:%s→%s⟧", checksumEdge(src, dst), src, dst)
}

func ParseEdges(data []byte) [][2]string {
	var result [][2]string
	re := regexp.MustCompile(`⟦([0-9a-fA-F]+):KoDeDgE:([a-zA-Z0-9\-.~_/]*?)→([a-zA-Z0-9\-.~_/]*?)⟧`)
	for _, m := range re.FindAllSubmatch(data, -1) {
		if len(m) != 4 {
			continue
		}
		sum, src, dst := string(m[1]), string(m[2]), string(m[3])
		if sum != checksumEdge(src, dst) {
			continue
		}
		result = append(result, [2]string{src, dst})
	}
	sort.Slice(result, func(i, j int) bool {
		if a, b := result[i][0], result[j][0]; a != b {
			return a < b
		}
		return result[i][1] < result[j][1]
	})
	return result
}

func checksumEdge(src, dst string) string {
	edge := fmt.Sprintf("KoDeDgE:%s→%s", src, dst)
	sum := sha256.Sum256([]byte(edge))
	return fmt.Sprintf("%0x", sum)[:8]
}
