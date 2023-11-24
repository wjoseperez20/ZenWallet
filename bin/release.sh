#!/usr/bin/env sh
current_branch=$(git branch --show-current)

if [ "$current_branch" != "main" ];then
    echo "The release command only can be run in main branch"
    echo "    (user \"git switch main\" to change to your main branch)"
    exit 10
fi

application_name="ZenWallet"
changelog_file="./CHANGELOG.md"
type_changes=("added" "changed" "fixed" "removed")
last_commit=$(git log --oneline --invert-grep --grep="docs: " --format="%h" -n1 )
date=$(date '+%Y-%m-%d')

last_revision=$(grep -Eo '\## \[v(.+?)\] .+' $changelog_file | head -1 | awk -F"[ ')]+" '/## \[/{print $5}' | sed 's/(//g' | sed 's/#//g') 

if [ $last_revision == $last_commit ]; then
    echo "There are not new commits to create a new version"
    exit 0
fi

for type in ${type_changes[@]}; do

    if [ "${type}" == "fixed" ];then
        fixed=$(git log --oneline --no-merges --format="%B" $last_revision..HEAD main | grep -E "^fix:|fix\(.*\):\s{1,}" | sed 's/fix:/*/g') 
    fi

    if [ "${type}" == "added" ];then
        added=$(git log --oneline --no-merges --format="%B" $last_revision..HEAD main | grep -E "^feat:|feat\(.*\):\s{1,}" | sed 's/feat:/*/g') 
    fi
done

if [ -z "$added" ] && [ -z "$fixed" ]; then
    echo "There are not new commits to create a new version"
    exit 0
fi

versions=$(grep -Eo '\## \[v(.+?)\]' $changelog_file | sed 's/##//g' | sed 's/\[//g' | sed 's/\]//g'  | sed 's/v//g')
last_version=$(echo $versions | sed 's/\ /\n/g' | head -1) 

IFS="." read -ra semantic_version <<< "$last_version"
echo "> current ${application_name} version --> [v${last_version}]\n"

major="${semantic_version[0]}"
minor="${semantic_version[1]}"
patch="${semantic_version[2]}"

echo "We are using semantic version for each new version, make you sure to create a proper version to deploy"
read -p "Enter the new version or enter for the suggested (${major}.${minor}.$((patch + 1))): " new_version

echo ""

if [ -z "$new_version" ]; then
  new_version="${major}.${minor}.$((patch + 1))"
fi

echo "the new version is: $new_version"

git checkout -b "docs/release-v$new_version"

insert_in_line=$(awk '/## \[Unreleased\]/{print NR}' $changelog_file)

{ head -n $insert_in_line $changelog_file; echo "\n## [v$new_version] - $date (#$last_commit)\n\n### [Added]\n\n$added\n\n### [Fixed]\n\n$fixed"; tail -n +"$((insert_in_line + 1))" $changelog_file; } > CHANGELOG.temp.md

rm $changelog_file
mv ./CHANGELOG.temp.md $changelog_file

git add $changelog_file
git commit -m "docs: new version $new_version added to chagelog file"
git push -o merge_request.create -o merge_request.target="main" -o merge_request.title="release v$new_version" origin "docs/release-v$new_version" 
