// 为保存按钮添加事件监听，拦截保存设置的 htmx 请求
document.getElementById('saveConfigButton').addEventListener('htmx:beforeRequest', function (event) {
    // 检查是否可以保存设置
    let selectedDir = document.getElementById('selectedDir').value;
    canSave = false;
    if (selectedDir === 'WorkingDirectory') {
        if (document.getElementById('ProgramDirectoryConfigDiv') === null && document.getElementById('HomeDirectoryConfigDiv') === null) {
            canSave = true;
        }
    }
    if (selectedDir === 'HomeDirectory') {
        if (document.getElementById('ProgramDirectoryConfigDiv') === null && document.getElementById('WorkingDirectoryConfigDiv') === null) {
            canSave = true;
        }
    }
    if (selectedDir === 'ProgramDirectory') {
        if (document.getElementById('HomeDirectoryConfigDiv') === null && document.getElementById('WorkingDirectoryConfigDiv') === null) {
            canSave = true;
        }
    }
    // 如果其他地方已经有配置了，则阻止请求并执行本地逻辑
    if (!canSave) {
        event.preventDefault();
        showToast('只允许在一个地方保存配置', 'info');
    }
});
document.getElementById('saveConfigButton').addEventListener('htmx:afterRequest', function (event) {
    // 只对ID为myButton的请求进行监听
    if (event.detail.successful) {
        showToast('保存设置文件成功！', 'info');
    } else {
        showToast('保存设置文件失败', 'error');
    }
});

// 为删除按钮添加事件监听，
document.getElementById('deleteConfigButton').addEventListener('htmx:beforeRequest', function (event) {
    // 检查是否可以删除配置
    let selectedDir = document.getElementById('selectedDir').value;
    canSelete = true;
    if (selectedDir === 'WorkingDirectory') {
        if (document.getElementById('WorkingDirectoryConfigDiv') === null) {
            canSelete = false;
        }
    }
    if (selectedDir === 'HomeDirectory') {
        if (document.getElementById('HomeDirectoryConfigDiv') === null) {
            canSelete = false;
        }
    }
    if (selectedDir === 'ProgramDirectory') {
        if (document.getElementById('ProgramDirectoryConfigDiv') === null) {
            canSelete = false;
        }
    }
    // 如果不满足条件，则阻止请求并执行本地逻辑
    if (!canSelete) {
        event.preventDefault();
        showToast('当前选择的路径下，没有可删除的配置文件', 'info');
    }
});
document.getElementById('deleteConfigButton').addEventListener('htmx:afterRequest', function (event) {
    // 只对ID为myButton的请求进行监听
    if (event.detail.successful) {
        showToast('删除设置文件成功！', 'info');
    } else {
        showToast('保删除设置文件失败', 'error');
    }
});