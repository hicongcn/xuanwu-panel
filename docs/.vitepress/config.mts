import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
    title: '玄武面板',
    description: '极致轻量、高性能的自动化任务调度平台',
    base: '/xuanwu-panel/',
    lang: 'zh-CN',
    themeConfig: {
        logo: '/logo.svg',
        nav: [
            { text: '快速开始', link: '/guide/introduction' },
            { text: '部署指南', link: '/guide/deployment' },
            { text: 'API 文档', link: '/guide/api' }
        ],

        sidebar: [
            {
                text: '基础指南',
                items: [
                    { text: '项目介绍', link: '/guide/introduction' },
                    { text: '部署说明', link: '/guide/deployment' },
                    { text: '开始使用', link: '/guide/getting-started' },
                    { text: 'API 文档', link: '/guide/api' }
                ]
            },
            {
                text: '功能指南',
                items: [
                    { text: '数据仪表', link: '/guide/dashboard' },
                    { text: '定时任务', link: '/guide/tasks' },
                    { text: '脚本管理', link: '/guide/scripts' },
                    { text: '执行历史', link: '/guide/history' },
                    { text: '变量机密', link: '/guide/environments' },
                    { text: '语言依赖', link: '/guide/languages' },
                    { text: '终端命令', link: '/guide/terminal' },
                    { text: '消息中心', link: '/guide/notify' },
                    { text: '仓库同步', link: '/guide/sync' },
                    { text: '命令行工具', link: '/guide/cli' },
                    {
                        text: '脚本示例',
                        link: '/guide/examples/',
                        items: [
                            { text: '浏览器示例', link: '/guide/examples/browser' },
                            { text: '内置库示例', link: '/guide/examples/builtin' }
                        ]
                    }
                ]
            },
            {
                text: '部署配置',
                items: [
                    { text: '系统配置', link: '/guide/configuration' },
                    { text: '反向代理', link: '/guide/nginx' }
                ]
            },
            {
                text: '其他',
                items: [
                    { text: '更新日志', link: '/guide/changelog' },
                    { text: '免责声明', link: '/guide/disclaimer' }
                ]
            }
        ],

        socialLinks: [
            { icon: 'github', link: 'https://github.com/hicongcn/xuanwu-panel' }
        ],

        footer: {
            message: 'Released under the Apache License 2.0.',
            copyright: 'Copyright © 2025-present hicongcn'
        },

        search: {
            provider: 'local'
        }
    },
    vite: {
        ssr: {
            noExternal: ['@scalar/api-reference']
        }
    }
})
