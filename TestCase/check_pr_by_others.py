import time 
import yaml
from pr_test_case import PullRequestOperation


with open('config.yaml', 'r') as f:
    info = yaml.load(f.read())['test case']
    owner = info[0]['owner']
    repo = info[1]['repo']
    local_owner = info[2]['local_owner']
    pr = PullRequestOperation(owner, repo, local_owner) 
   
    print('Prepare:')
    print('step 1/4: git clone')
    pr.git_clone()
    print('\nstep 2/4: change file')
    pr.change_file()
    print('\nstep 3/4: git push')
    pr.git_push()
    print('\nstep 4/4: pull request')
    number = pr.pull_request()
    print('the number of the pull request: {}'.format(number))
    time.sleep(10)

    print('\n\nTest:')
    print('test case: check_pr_by_others')
    labels = pr.get_all_labels(number)
    print('labels: {}'.format(labels))
    print('add labels to pr:')
    pr.add_labels_2_pr(number, '["lgtm", "approved"]')
    labels = pr.get_all_labels(number)
    print('labels: {}'.format(labels))
    if 'lgtm' in labels and 'approved' in labels and 'ci-bot-cla/yes' in labels:
        pr.comment_by_others(number, '/check-pr')
        time.sleep(10)
        code = pr.get_pr_status(number)
        if code == 200:
            print('test case check_pr_by_others succeeded')
        else:
            print('failed code: {}'.format(code))
            print('test case check_pr_by_others failed') 
    else:
        if 'lgtm' not in labels:
            print('need "lgtm" to merge')
        if 'approved' not in labels:
            print('need "approved" to merge')
        if 'ci-bot-cla/yes' not in labels:
            print('need "ci-bot-cla/yes" to merge')
        print('test case check_pr_by_others failed')
