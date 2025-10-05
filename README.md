# qbittorrent-tool

一个用于管理 qBittorrent 的工具，支持自动分类、标签设置、做种限制等功能。

## 功能特性

- 根据 tracker 域名自动设置标签（如 baidu.com）
- 根据种子保存路径自动设置分类
- 基于标签、分类、tracker、关键词等条件精确控制做种时间、分享率、上传速度，并支持删种或例外保种
- 根据 tracker 工作状态（尤其是错误状态）自动添加标签以便识别问题
- 批量替换 tracker（待实现）

## 安装

```bash
# 克隆项目
git clone https://github.com/fengqi/qbittorrent-tool.git
cd qbittorrent-tool

# 安装依赖
pip install -r requirements.txt
```

## 使用方法

```bash
# 直接运行
python main.py

# 或者安装后运行
pip install -e .
qbittorrent-tool
```

### 打包为可执行文件

```bash
# 安装打包工具
pip install pyinstaller

# 打包
make build

# 或者直接使用 pyinstaller
pyinstaller pyinstaller.spec
```

打包后的可执行文件位于 `dist/` 目录中。

## 配置文件

请参考 `example.config.json` 创建你的配置文件，并使用 `-c` 参数指定配置文件路径：

```bash
python main.py -c ./config.json
```

## 技术栈

Python 3.6+

## 许可证

MIT
