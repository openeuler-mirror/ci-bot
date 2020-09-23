## How to run the test case scripts  
### prepare  
- config owner and repo  
    
  The owner and repo were default.Config your local owner to specify where you forked to.  
  ```
  - owner: ci-bot
  - repo: community
  - local_owner: ******
  ``` 
 - config access_token  
  
   To run the test case scripts, maybe you'll prepare one more gitee access_token because most gitee apis need   
  access_token.  
    
### start  
  If you see the following,the script works.  
``` 
Prepare:
step 1/4: git clone
Cloning into 'community'...
remote: Enumerating objects: 355, done.
remote: Counting objects: 100% (355/355), done.
remote: Compressing objects: 100% (312/312), done.
remote: Total 355 (delta 171), reused 154 (delta 23), pack-reused 0
Receiving objects: 100% (355/355), 1.05 MiB | 440.00 KiB/s, done.
Resolving deltas: 100% (171/171), done.

step 2/4: change file
remove test.txt

step 3/4: git push
[master 5aa1daa] test
 1 file changed, 0 insertions(+), 0 deletions(-)
 delete mode 100644 test.txt
Counting objects: 2, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (1/1), done.
Writing objects: 100% (2/2), 212 bytes | 212.00 KiB/s, done.
Total 2 (delta 1), reused 1 (delta 1)
remote: Powered by GITEE.COM [GNK-5.0]
To https://gitee.com/******/community.git
   7d62a36..5aa1daa  master -> master

step 4/4: pull request
the number of the pull request: ******


Test:
test case 1: without comments by contributor
labels: ['ci-bot-cla/yes']
test case 1 succeeded

test case 2: /lgtm
comment body: /lgtm
labels: ['ci-bot-cla/yes']
test case 2 succeeded

test case 3: comment /lgtm by others
comment body: /lgtm
labels: ['ci-bot-cla/yes']
test case 3 succeeded

test case 4: comment /approve by others
comment body: /approve
labels: ['ci-bot-cla/yes']
test case 4 succeeded

test case 5: /approve
comment body: /approve
labels: ['approved', 'ci-bot-cla/yes']
test case 5 succeeded

test case 6: tag stat/need-squash
labels_before_commit: ['approved', 'ci-bot-cla/yes']
[master 19fce3f] change test.txt
 1 file changed, 1 insertion(+)
 create mode 100644 test.txt
Counting objects: 3, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (2/2), done.
Writing objects: 100% (3/3), 299 bytes | 299.00 KiB/s, done.
Total 3 (delta 1), reused 1 (delta 0)
remote: Powered by GITEE.COM [GNK-5.0]
To https://gitee.com/******/community.git
   5aa1daa..19fce3f  master -> master
lables_after_commit: ['approved', 'ci-bot-cla/yes', 'stat/need-squash']
test case 6 succeeded

test case 7: add labels
labels: ['lgtm', 'approved', 'ci-bot-cla/yes', 'stat/need-squash']
test case 7 succeeded

test case 8: check-pr
comment body: /check-pr
test case 8 succeeded
```
### notice
- You can tag `lgtm` to the pull request if you're the repository member, but it doesn't work when commenting `/lgtm` under the pull request you created.
- If you want to merge a pull request, you must have permissions to tag `approve` to the pull request.
- If you can edit a pull request, you can also invoke the gitee apis to create labels or delete labels.
- The pull request may be merged when `ci-bot-cla/yes`, `lgtm` and `approved` both in the pr labels.Anyone can comment `/check-pr` to trigger the merge action if it's still open.
