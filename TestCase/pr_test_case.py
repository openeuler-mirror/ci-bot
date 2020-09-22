import os
import requests
import subprocess
import sys
import time


class PullRequestOperation(object):
    def __init__(self, owner, repo, local_owner):
        """initialize owner, repo and access_token"""
        self.owner = owner
        self.repo = repo
        self.local_owner = local_owner
        self.access_token = os.getenv('ACCESS_TOKEN', '')

    def git_clone(self):
        """git clone code"""
        subprocess.call("git clone https://gitee.com/{}/{}.git".format(self.local_owner, self.repo), shell=True)

    def change_file(self):
        """change file: Test whether test.txt exists.Remove test.txt if it exists, or touch test.txt"""
        subprocess.call(
            "cd {}/; [ -f test.txt ]; if [ $? -eq 0 ]; then rm test.txt; echo 'remove test.txt'; else touch test.txt; echo 'touch test.txt'; fi".format(self.repo), shell=True)

    def write_2_file(self):
        """write some info to the test file"""
        subprocess.call(
            'cd {}/; echo "hello" > test.txt; git add .; git commit -m "change test.txt"; git push'.format(self.repo),
            shell=True)

    def git_push(self):
        """push code"""
        subprocess.call("cd {}/; git add . ; git commit -m 'test'; git push".format(self.repo), shell=True)

    def pull_request(self):
        """create a pull request"""
        head = '{}:master'.format(self.local_owner)
        data = {
            'access_token': self.access_token,
            'title': 'test',
            'head': head,
            'base': 'master'
        }
        url = 'https://gitee.com/api/v5/repos/{}/{}/pulls'.format(self.owner, self.repo)
        r = requests.post(url, data)
        if r.status_code == 201:
            number = r.json()['number']
            return number
        if r.status_code == 400:
            number = r.json()['message'].split('!')[1].split(' ')[0]
            return number

    def comment(self, number, body):
        """comment under the pull request"""
        data = {
            'access_token': self.access_token,
            'body': body
        }
        url = 'https://gitee.com/api/v5/repos/{}/{}/pulls/{}/comments'.format(self.owner, self.repo, number)
        print('comment body: {}'.format(body))
        requests.post(url, data)

    def comment_by_others(self, number, body):
        """comment under the pull request"""
        data = {
            'access_token': os.getenv('ACCESS_TOKEN_TWO', ''),
            'body': body
        }
        url = 'https://gitee.com/api/v5/repos/{}/{}/pulls/{}/comments'.format(self.owner, self.repo, number)
        print('comment body: {}'.format(body))
        requests.post(url, data)

    def get_all_comments(self, number):
        """get all comments under the pull request"""
        params = 'access_token={}'.format(self.access_token)
        url = 'https://gitee.com/api/v5/repos/{}/{}/pulls/{}/comments?per_page=100'.format(self.owner, self.repo,
                                                                                           number)
        r = requests.get(url, params)
        comments = []
        if r.status_code == 404:
            return r.json()['message']
        if r.status_code == 200:
            if len(r.json()) > 0:
                for comment in r.json():
                    comments.append(comment['body'])
                return comments
            else:
                return comments

    def get_all_labels(self, number):
        """get all labels belong to the pull request"""
        params = 'access_token={}'.format(self.access_token)
        url = 'https://gitee.com/api/v5/repos/{}/{}/pulls/{}/labels'.format(self.owner, self.repo, number)
        r = requests.get(url, params)
        labels = []
        if r.status_code == 200:
            if len(r.json()) > 0:
                for i in r.json():
                    labels.append(i['name'])
                return labels
            else:
                return labels


if __name__ == '__main__':
    try:
        owner = sys.argv[1]
        repo = sys.argv[2]
        local_owner = sys.argv[3]
        pr = PullRequestOperation(owner, repo, local_owner)

        print('step 1: git clone')
        pr.git_clone()
        print('\nstep 2: change file')
        pr.change_file()
        print('\nstep 3: git push')
        pr.git_push()
        print('\nstep 4: pull request')
        number = pr.pull_request()
        print('the number of the pull request: {}'.format(number))
        time.sleep(10)

        print('\ntest case 1: without comments by contributor')
        comments = pr.get_all_comments(number)
        labels = pr.get_all_labels(number)
        print('labels: {}'.format(labels))
        errors = 0
        if len(comments) == 0:
            print('no "Welcome to ci-bot Community."')
            print('no "Thanks for your pull request."')
        else:
            if 'Welcome to ci-bot Community.' not in comments[0]:
                print('no "Welcome to ci-bot Community."')
                errors += 1
            if 'Thanks for your pull request.' not in comments[-1]:
                print('no "Thanks for your pull request."')
                errors += 1
        if len(labels) == 0:
            print('no label "ci-bot-cla/yes" or "ci-bot-cla/no"')
            errors += 1
        elif len(labels) > 0:
            if 'ci-bot-cla/yes' not in labels:
                print('no label "ci-bot-cla/yes"')
                errors += 1
                if 'ci-bot-cla/no' not in labels:
                    print('no label "ci-bot-cla/no"')
                    errors += 1
        if errors == 0:
            print('test case 1 succeeded.')

        print('\ntest case 2: /lgtm')
        pr.comment(number, '/lgtm')
        time.sleep(10)
        labels = pr.get_all_labels(number)
        print('labels: {}'.format(labels))
        comments = pr.get_all_comments(number)
        if 'can not be added in your self-own pull request' in comments[-1]:
            print('test case 2 succeeded.')
        else:
            print('test case 2 failed.')
            print(comments[-1])

        print('\ntest case 3: comment /lgtm by others')
        pr.comment_by_others(number, '/lgtm')
        time.sleep(10)
        labels = pr.get_all_labels(number)
        print('labels: {}'.format(labels))
        comments = pr.get_all_comments(number)
        if 'Thanks for your review' in comments[-1]:
            print('test case 3 succeeded')
        else:
            print('test case 3 failed')
            print(comments[-1])

        print('\ntest case 4: /approve')
        pr.comment(number, '/approve')
        time.sleep(10)
        labels = pr.get_all_labels(number)
        print('labels: {}'.format(labels))
        comments = pr.get_all_comments(number)
        if '***approved*** is added in this pull request by' in comments[-1]:
            print('test case 4 succeeded.')
        else:
            print('test case 4 failed.')
            print(comments[-1])

        print('\ntest case 5: comment /approve by others')
        pr.comment_by_others(number, '/approve')
        time.sleep(10)
        labels = pr.get_all_labels(number)
        print('labels: {}'.format(labels))
        comments = pr.get_all_comments(number)
        if 'has no permission to add' in comments[-1]:
            print('test case 5 succeeded')
        else:
            print('test case 5 failed')
            print(comments[-1])

        print('\ntest case 6: tag stat/need-squash')
        labels_before_commit = pr.get_all_labels(number)
        print('labels_before_commit: {}'.format(labels_before_commit))
        pr.write_2_file()
        time.sleep(10)
        lables_after_commit = pr.get_all_labels(number)
        print('lables_after_commit: {}'.format(lables_after_commit))
        if 'lgtm' not in labels and 'stat/need-squash' in lables_after_commit:
            print('test case 6 succeeded')
        else:
            print('test case 6 failed')
    except IndexError:
        print('\n3 arguments were needed, please check!\n')
