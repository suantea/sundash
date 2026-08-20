/**
 * Built-in navigation templates for the "Import template" feature.
 * Each template is a list of groups + cards that can be appended to the
 * current panel via doImport().
 */

export interface TemplateCard {
  title: string
  url: string
  icon?: string
  description?: string
}

export interface TemplateGroup {
  name: string
  cards: TemplateCard[]
}

export interface ImportTemplate {
  name: string
  groups: TemplateGroup[]
}

export const importTemplates: ImportTemplate[] = [
  {
    name: '常用导航',
    groups: [
      {
        name: '搜索引擎',
        cards: [
          { title: 'Google', url: 'https://www.google.com', icon: 'logos:google-icon', description: '全球最大搜索引擎' },
          { title: 'Bing', url: 'https://www.bing.com', icon: 'logos:bing', description: '微软搜索引擎' },
          { title: '百度', url: 'https://www.baidu.com', icon: 'logos:baidu', description: '国内最大搜索引擎' },
          { title: 'DuckDuckGo', url: 'https://duckduckgo.com', icon: 'logos:duckduckgo-icon', description: '不追踪的搜索引擎' },
          { title: 'Yandex', url: 'https://yandex.com', icon: 'logos:yandex-icon', description: '俄罗斯搜索引擎' },
        ],
      },
      {
        name: '开发工具',
        cards: [
          { title: 'GitHub', url: 'https://github.com', icon: 'logos:github-icon', description: '全球最大代码托管平台' },
          { title: 'GitLab', url: 'https://gitlab.com', icon: 'logos:gitlab', description: '开源代码托管平台' },
          { title: 'Stack Overflow', url: 'https://stackoverflow.com', icon: 'logos:stackoverflow-icon', description: '程序员问答社区' },
          { title: 'MDN Web Docs', url: 'https://developer.mozilla.org', icon: 'logos:mdn-web-docs', description: 'Web 开发者文档' },
          { title: 'Docker Hub', url: 'https://hub.docker.com', icon: 'logos:docker-icon', description: 'Docker 镜像仓库' },
          { title: 'npm', url: 'https://www.npmjs.com', icon: 'logos:npm-icon', description: 'Node 包管理仓库' },
        ],
      },
      {
        name: '社交媒体',
        cards: [
          { title: 'Twitter', url: 'https://twitter.com', icon: 'logos:twitter', description: '全球社交平台' },
          { title: 'YouTube', url: 'https://www.youtube.com', icon: 'logos:youtube-icon', description: '全球最大视频平台' },
          { title: 'Reddit', url: 'https://www.reddit.com', icon: 'logos:reddit-icon', description: '兴趣社区' },
          { title: '微博', url: 'https://weibo.com', icon: 'logos:sina-weibo', description: '国内社交平台' },
          { title: 'Bilibili', url: 'https://www.bilibili.com', icon: 'logos:bilibili-icon', description: '国内视频平台' },
        ],
      },
      {
        name: '云服务',
        cards: [
          { title: 'AWS Console', url: 'https://console.aws.amazon.com', icon: 'logos:aws', description: '亚马逊云控制台' },
          { title: 'Google Cloud', url: 'https://console.cloud.google.com', icon: 'logos:google-cloud', description: '谷歌云控制台' },
          { title: 'Azure', url: 'https://portal.azure.com', icon: 'logos:azure-icon', description: '微软云控制台' },
          { title: 'Cloudflare', url: 'https://dash.cloudflare.com', icon: 'logos:cloudflare', description: 'CDN 与域名管理' },
          { title: 'Vercel', url: 'https://vercel.com/dashboard', icon: 'logos:vercel-icon', description: '前端部署平台' },
          { title: '阿里云', url: 'https://home.console.aliyun.com', icon: 'logos:alibaba-cloud', description: '阿里云控制台' },
        ],
      },
      {
        name: '实用工具',
        cards: [
          { title: 'Wikipedia', url: 'https://www.wikipedia.org', icon: 'logos:wikipedia-icon', description: '维基百科' },
          { title: 'Notion', url: 'https://www.notion.so', icon: 'logos:notion-icon', description: '笔记与协作工具' },
          { title: 'Figma', url: 'https://www.figma.com', icon: 'logos:figma', description: '在线设计工具' },
          { title: 'Speedtest', url: 'https://www.speedtest.net', icon: 'logos:ookla-speedtest', description: '网络测速' },
          { title: 'IP 查询', url: 'https://ip.sb', icon: 'mdi:ip-network', description: '查询本机 IP' },
          { title: 'Unicode 表', url: 'https://unicode-table.com', icon: 'mdi:format-text-variant', description: 'Unicode 字符查询' },
        ],
      },
    ],
  },
]
