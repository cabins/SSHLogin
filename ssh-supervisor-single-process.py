# ！/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import time

import paramiko


def ssh_run(ip='127.0.0.1', port=22, username='root', password='', cmd=''):
    """
    登陆到ssh服务器，并进行相应的操作处理
    """
    print('Starts running on ', ip)
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy)

    # 登陆远程机器
    try:
        ssh.connect(ip, int(port), username, password)
    except paramiko.ssh_exception.AuthenticationException as e:
        # 错误信息写入文件
        error_write(ip, '用户名或密码错误')
        ssh.connect(ip, int(port), "NEW_USERNAME", "NEW_PASSWORD")
        return
    except TimeoutError as e:
        # 错误信息写入文件
        error_write(ip, '远程主机没有反应，连接失败')
        return

    stdin, stdout, stderr = ssh.exec_command(cmd)

    err = stderr.read().decode('utf-8')
    out = stdout.read().decode('utf-8')
    # 如果执行有错误 #
    print(err)
    if err != '':
        # 错误信息写入文件
        error_write(ip, err)
        if out != '':
            content_write(ip, out)
    else:
        # 返回结果信息写入文件
        content_write(ip, out)

    # 断开连接
    ssh.close()


def error_write(ip, message):
    """
    错误信息的保存
    """
    # 注意：如果要用t在文件路径中的话，Windows上会报错，需要把其中的：替换成其他非路径分隔符，如_
    t = time.strftime('%Y-%m-%d %H:%M:%S')

    with open(os.path.join('error', '.'.join(['error', ip, 'txt'])), 'a', encoding='utf-8') as f:
        f.write('\t'.join([t, message]) + '\n')


def content_write(ip, message):
    """
    正常结果的保存
    """
    t = time.strftime('%Y-%m-%d %H:%M:%S')

    with open(os.path.join('result', '.'.join(['result', ip, 'txt'])), 'a', encoding='utf-8') as f:
        f.write('\n'.join([t, message]) + '\n')


def main():
    """从ips.csv中获取IP，登录信息，需要执行的命令等
       ips.csv符合csv的格式，逗号隔开，没有标题栏，每列代表IP, 端口号，用户名，密码，需要执行的命令
       需要执行的命令用分号隔开
    """

    if not (os.path.exists('result') and os.path.isdir('result')):
        os.mkdir('result')

    if not (os.path.exists('error') and os.path.isdir('error')):
        os.mkdir('error')

    with open('ips.csv') as f:
        for line in f.readlines():
            ip, port, username, password, cmd = line.split(',')

            ssh_run(ip, int(port), username, password, cmd)


if __name__ == '__main__':
    main()
