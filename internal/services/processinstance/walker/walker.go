package walker

import (
	"context"
	"fmt"

	d "github.com/grafvonb/camunder/internal/domain"
	"github.com/grafvonb/camunder/internal/services"
	"github.com/grafvonb/camunder/internal/services/processinstance"
)

func Ancestry(ctx context.Context, s processinstance.API, startKey int64) (rootKey int64, path []int64, chain map[int64]d.ProcessInstance, err error) {
	// visited keeps track of visited nodes to detect cycles
	// well-know pattern to have fast lookups, no duplicates, clear semantic and low memory usage with visited[cur] = struct{}{} below
	visited := make(map[int64]struct{})
	chain = make(map[int64]d.ProcessInstance)

	cur := startKey
	for {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return 0, nil, chain, ctx.Err()
		default:
		}

		if _, seen := visited[cur]; seen {
			return 0, nil, chain, fmt.Errorf("%w for this key %d", services.ErrCycleDetected, cur)
		}
		visited[cur] = struct{}{}

		it, getErr := s.GetProcessInstanceByKey(ctx, cur)
		if getErr != nil {
			return 0, nil, chain, fmt.Errorf("get %d: %w", cur, getErr)
		}
		chain[cur] = it
		path = append(path, cur)

		// no parent => cur is root
		if it.ParentKey == 0 {
			rootKey = cur
			return
		}

		cur = it.ParentKey
	}
}

func Descendants(ctx context.Context, s processinstance.API, rootKey int64) (desc []int64, edges map[int64][]int64, chain map[int64]d.ProcessInstance, err error) {
	visited := make(map[int64]struct{})
	edges = make(map[int64][]int64)
	chain = make(map[int64]d.ProcessInstance)

	// depth-first search (DFS) to explore the tree
	var dfs func(int64) error
	dfs = func(parent int64) error {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if _, seen := visited[parent]; seen {
			// already expanded this subtree
			return nil
		}
		visited[parent] = struct{}{}

		desc = append(desc, parent)
		if _, ok := chain[parent]; !ok {
			it, getErr := s.GetProcessInstanceByKey(ctx, parent)
			if getErr != nil {
				return fmt.Errorf("get %d: %w", parent, getErr)
			}
			chain[parent] = it
		}

		children, e := s.GetDirectChildrenOfProcessInstance(ctx, parent)
		if e != nil {
			return fmt.Errorf("list children of %d: %w", parent, e)
		}

		// keep an entry even if no children (useful for tree rendering)
		if _, ok := edges[parent]; !ok {
			edges[parent] = nil
		}

		for i := range children {
			it := children[i]
			k := it.Key

			edges[parent] = append(edges[parent], k)
			chain[k] = it

			if dfsErr := dfs(k); dfsErr != nil {
				return dfsErr
			}
		}
		return nil
	}

	if err = dfs(rootKey); err != nil {
		return nil, nil, nil, err
	}
	return desc, edges, chain, nil
}

func Family(ctx context.Context, s processinstance.API, startKey int64) (fam []int64, edges map[int64][]int64, chain map[int64]d.ProcessInstance, err error) {
	rootKey, _, _, err := Ancestry(ctx, s, startKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ancestry fetch: %w", err)
	}
	return Descendants(ctx, s, rootKey)
}
