PHONY: git/prune
# Prunes all branches that have been deleted on the remote & tags that are no longer reachable.
git/prune:
	@git fetch -pP && for branch in $$(git for-each-ref --format '%(refname) %(upstream:track)' refs/heads | awk '$$2 == "[gone]" {sub("refs/heads/", "", $$1); print $$1}'); do git branch -D $$branch; done && echo "Branch cleanup complete." || echo "Branch cleanup failed or nothing to prune."