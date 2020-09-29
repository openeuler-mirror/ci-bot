import subprocess
import time
import yaml
from pr_test_case import PullRequestOperation


with open('config.yaml', 'r') as f:
    info = yaml.load(f.read())['test case']
    owner = info[0]['owner']
    repo = info[1]['repo']
    local_owner = info[2]['local_owner']
    pr = PullRequestOperation(owner, repo, local_owner)
    # create a SIG with as least 1 repo
    print('step 1/4: git clone')
    pr.git_clone()
    print('\nstep 2/4: change file')
    # new a sig dir contain OWNERS
    subprocess.call("cd community/sig; mkdir Container; cd Container; wget https://gitee.com/openeuler/community/raw/master/sig/Container/OWNERS", shell=True)
    # add repo to sig/sigs.yaml
    with open('community/sig/sigs.yaml', 'a') as f2:
        f2.write('\n- name: Container\n  repositories:\n  - ci-bot/kubernetes')
    # add repo description to repository/gerogecaotest.yaml
    with open('community/repository/georgecaotest.yaml', 'a') as f3:
        f3.write('\n- name: kubernetes\n  description: "kubernetes employ and operate"\n  protected_branches:\n  - master\n  type: public')
    print('\nstep 3/4: git push')
    pr.git_push()
    print('\nstep 4/4: pull request')
    number = pr.pull_request()
    print('the number of the pull request: {}'.format(number))
    time.sleep(10)
    # tag 'lgtm' and 'approved' to merge
    pr.add_labels_2_pr(number, '["lgtm"]')
    # If you're sure the pr has no trouble and you wish the pr to be merged automatically, then you can turn off the comment below.
    #pr.comment(number, '/approve')
