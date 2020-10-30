import requests
import yaml


def compare_repo_members_with_maintainers(sig_name):
    maintainers = []
    with open('community/sig/{}/OWNERS'.format(sig_name), 'r') as f:
        res = yaml.load(f.read())
        for i in res['maintainers']:
            maintainers.append(i)
    print('maintainers of {}: {}'.format(sig_name, sorted(maintainers)))
    with open('community/sig/sigs.yaml') as f2:
        sigs = yaml.load(f2.read())['sigs']
        for sig in sigs:
            if sig['name'] == sig_name:
                for repo in sig['repositories']:
                    repo_members = []
                    url = 'https://gitee.com/api/v5/repos/{}/collaborators'.format(repo)
                    r = requests.get(url)
                    for member in r.json():
                        repo_members.append(member['login'])
                    repo_members.remove('georgecao')
                    # compare
                    print('repo_members of {}: {}'.format(repo, sorted(repo_members)))
                    if sorted(repo_members) == sorted(maintainers):
                        print('success')
                    else:
                        print('the repo members should be equal to maintainers')
