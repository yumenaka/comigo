#!/usr/bin/env python3
"""用伪终端和本地下载替身验证选择、安装及校验，不联网、不写系统目录。"""
import hashlib
import io
import os
from pathlib import Path
import pty
import select
import subprocess
import tarfile
import tempfile
import time

SCRIPT = Path(__file__).resolve().parents[1] / "get.sh"


def run(env, args, answers=None):
    """等待终端输出读完再回收进程，避免遗漏退出前的日志。"""
    command = ["bash", str(SCRIPT), "--version", "v9.9.9", *args]
    if answers is None:
        result = subprocess.run(command, env=env, start_new_session=True,
                                stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                                stderr=subprocess.STDOUT, timeout=10)
        return result.returncode, result.stdout.decode()
    pid, fd = pty.fork()
    if pid == 0:
        os.execvpe("bash", command, env)
    output = b""
    pending = list(answers)
    deadline = time.monotonic() + 10
    status = None
    try:
        while time.monotonic() < deadline:
            if select.select([fd], [], [], 0.1)[0]:
                try:
                    chunk = os.read(fd, 65536)
                except OSError:
                    break
                if not chunk:
                    break
                output += chunk
                while pending and pending[0][0].encode() in output:
                    prompt, answer = pending.pop(0)
                    os.write(fd, (answer + "\n").encode())
        else:
            raise AssertionError("Installer timed out: " + output.decode())
        assert not pending, (pending, output.decode())
        if pid is not None:
            _, status = os.waitpid(pid, 0)
            pid = None
        return os.waitstatus_to_exitcode(status), output.decode()
    finally:
        os.close(fd)
        if pid is not None:
            done, _ = os.waitpid(pid, os.WNOHANG)
            if not done:
                os.kill(pid, 9)
                os.waitpid(pid, 0)


with tempfile.TemporaryDirectory(prefix="comigo-installer-test-") as tmp:
    home = Path(tmp)
    mock = home / "mock"
    mock.mkdir()
    for name, content in {
        "curl": '#!/bin/bash\nfor url do :; done\nprintf "DOWNLOAD %s\\n" "$url"\nexit 1\n',
        "uname": '#!/bin/bash\ncase "$1" in -s) echo "${TEST_OS:-Linux}";; -m) echo x86_64;; esac\n',
    }.items():
        p = mock / name
        p.write_text(content)
        p.chmod(0o755)
    env = {**os.environ, "HOME": tmp, "PATH": f"{mock}:{home}/.local/bin:/usr/local/bin:/usr/bin:/bin",
           "LC_ALL": "C", "COMIGO_INSTALL_DIR": ""}
    for source, directory, path in [
        ("1", "1", "/usr/bin"),
        ("2", "2", "/usr/local/bin"),
        ("1", "3", str(home / ".local/bin")),
    ]:
        _, output = run(env, [], [("Enter 1 or 2:", source), ("Enter a number", directory)])
        assert "destination: " + path + "/comi" in output, output
        base = "https://github.com/" if source == "1" else "https://comigo.xyz/yumenaka/"
        assert "DOWNLOAD " + base in output, output
        assert output.index("1) /usr/bin") < output.index("2) /usr/local/bin") < output.index("3) " + str(home / ".local/bin")), output
        assert "/usr/bin (root required" in output and str(home / ".local/bin") + " (no root required)" in output, output
        assert str(home / "bin") not in output, output
        assert "Download failed" in output, output
    # 空输入与无效编号不能触发默认选择。
    _, output = run(env, [], [("Enter 1 or 2:", "\ninvalid\n2"),
                              ("Enter a number", "\n9\n3")])
    assert "DOWNLOAD https://comigo.xyz/" in output, output
    assert "destination: " + str(home / ".local/bin/comi") in output, output
    for args, expected in [
        ([], "Specify --github or --cn"),
        (["--github"], "Specify --install-dir or --system"),
        (["--github", "--install-dir", str(home / "custom path")], "DOWNLOAD https://github.com/"),
        (["--cn", "--system"], "DOWNLOAD https://comigo.xyz/"),
    ]:
        code, output = run(env, args)
        assert code != 0 and expected in output, output
    _, output = run({**env, "TEST_OS": "Darwin"}, ["--arch", "arm64"],
                    [("Enter 1 or 2:", "1"), ("Enter a number", "1")])
    assert "1) /usr/local/bin" in output and ") /usr/bin" not in output and "MacOS/arm64" in output, output
    # PATH 过滤不受 PATH 顺序、重复项或尾部斜杠影响；不匹配相似目录名。
    for path, expected in [
        (f"{mock}:/bin:/usr/local/bin/", "/usr/local/bin"),
        (f"{mock}:{home}/.local/bin:/bin", str(home / ".local/bin")),
        (f"{mock}:/usr/local/bin:/usr/bin:/usr/bin:/bin", "/usr/bin"),
    ]:
        code, output = run({**env, "PATH": path}, ["--github"], [("Enter a number", "1")])
        assert code != 0 and "destination: " + expected + "/comi" in output, output
        assert "1) " + expected in output, output
        if expected == "/usr/local/bin":
            assert ") /usr/bin" not in output and ") " + str(home / ".local/bin") not in output, output
    code, output = run({**env, "PATH": f"{mock}:/bin:/usr/bin-other"}, ["--github"], [])
    assert code != 0 and "No supported install directory is in PATH" in output, output
    assert "DOWNLOAD" not in output, output
    assert not (home / ".local/bin").exists()
    assert not (home / "bin").exists()

    # 下载替身只读取临时发布包；实际解包、校验与原子替换仍走安装器。
    assets = home / "assets"
    assets.mkdir()
    archive = assets / "comi_v9.9.9_Linux_x86_64.tar.gz"
    binary = b"#!/bin/bash\necho 'Comigo v9.9.9'\n"
    with tarfile.open(archive, "w:gz") as bundle:
        entry = tarfile.TarInfo("comi")
        entry.size = len(binary)
        entry.mode = 0o755
        bundle.addfile(entry, io.BytesIO(binary))
    checksums = assets / "checksums.txt"
    checksums.write_text(hashlib.sha256(archive.read_bytes()).hexdigest() + "  " + archive.name + "\n")
    (mock / "curl").write_text('''#!/bin/bash
while [[ $# -gt 0 ]]; do
    case "$1" in --output) output=$2; shift;; esac
    url=$1; shift
done
printf 'DOWNLOAD %s\n' "$url"
cp -- "$TEST_ASSETS/${url##*/}" "$output"
''')
    env["TEST_ASSETS"] = str(assets)
    target = home / "custom path" / "comi"
    for source in ("--github", "--cn"):
        code, output = run(env, [source, "--force", "--install-dir", str(target.parent)])
        assert code == 0 and target.read_bytes() == binary, output
        assert target.stat().st_mode & 0o111, output
    original = b"#!/bin/bash\necho 'Comigo v1.0.0'\n"
    target.write_bytes(original)
    checksums.write_text("0" * 64 + "  " + archive.name + "\n")
    code, output = run(env, ["--github", "--force", "--install-dir", str(target.parent)])
    assert code != 0 and "SHA-256 mismatch" in output, output
    assert target.read_bytes() == original
    assert not list(target.parent.glob(".comigo-install.*"))
print("Installer choices, successful installation and checksum protection passed")
