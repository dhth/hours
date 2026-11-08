#!/usr/bin/env bash
# Posts screenshots from capture.sh to a pull request as a comment, with the
# images attached to the comment itself.
#
# usage: publish.sh <screenshots-dir> <pr-number> [commit]
#
# commit is the commit the screenshots were taken from; it defaults to HEAD.
# The comment links it on the repo this checkout's origin points at. To post
# to a PR in another repo, set GH_REPO=<owner>/<repo>.
# Needs gh 2.99.0 or newer (for --attach), signed in with push access.
set -euo pipefail

if [ $# -lt 2 ] || [ $# -gt 3 ]; then
	echo "usage: $0 <screenshots-dir> <pr-number> [commit]" >&2
	exit 2
fi

DIR=$1
PR=$2
SHA=$(git rev-parse "${3:-HEAD}")

# Match on the saved text: piping into grep -q fails under pipefail when grep
# exits early and gh gets SIGPIPE.
help=$(gh pr comment --help 2>/dev/null || true)
case $help in *--attach*) ;; *)
	echo "error: gh $(gh --version | head -n 1 | cut -d' ' -f3) can't attach images; upgrade to 2.99.0 or newer" >&2
	exit 1
	;;
esac
[ -f "$DIR/scenarios.tsv" ] || {
	echo "error: $DIR/scenarios.tsv not found; run capture.sh first" >&2
	exit 1
}

SOURCE=$(git remote get-url origin | sed -E 's#^(https://github\.com/|git@github\.com:)##; s#\.git$##')
body="$DIR/comment.md" # image paths in the body are relative to this file
attach=()
failed=

{
	echo "### Screenshots"
	echo
	echo "Taken from [\`${SHA:0:7}\`](https://github.com/$SOURCE/commit/$SHA)."
	echo
	# collapsed by default; GitHub needs the blank line after </summary>
	echo "<details><summary>Screenshots ($(wc -l <"$DIR/scenarios.tsv" | tr -d ' '))</summary>"
	echo
	echo "| Scenario | Screenshot |"
	echo "| --- | --- |"
	while IFS="$(printf '\t')" read -r name desc; do
		if [ -f "$DIR/$name.png" ]; then
			echo "| $desc | ![$name](./$name.png) |"
			attach+=(--attach "$DIR/$name.png")
		elif [ -f "$DIR/FAILED-$name.png" ]; then
			echo "| ❌ **Failed:** $desc | ![FAILED-$name](./FAILED-$name.png) |"
			attach+=(--attach "$DIR/FAILED-$name.png")
			failed=$name
		else # the step failed without a screen to save
			echo "| ❌ **Failed:** $desc | (no screen captured) |"
			failed=$name
		fi
	done <"$DIR/scenarios.tsv"
	echo
	echo "</details>"

	if [ -n "$failed" ] && [ -f "$DIR/FAILED-$failed.txt" ]; then
		echo
		echo "<details><summary>Screen text when \"$failed\" failed</summary>"
		echo
		echo '```text'
		cat "$DIR/FAILED-$failed.txt"
		echo # termctrl's .txt has no trailing newline; the fence needs its own line
		echo '```'
		echo
		echo "</details>"
	fi
} >"$body"

gh pr comment "$PR" --body-file "$body" ${attach[@]+"${attach[@]}"}
