import { createRequire } from 'module'
import { defineConfig, type DefaultTheme } from 'vitepress'

const require = createRequire(import.meta.url)
const pkg = require('../../package.json')

export const en = defineConfig({
  lang: 'en-US',
  description: 'A CLI tool for building Go applications.',

  themeConfig: {
    nav: nav(),

    sidebar: {
      '/en/': { base: '/en/', items: sidebarGuide() },
    },

    editLink: {
      pattern: 'https://github.com/go-nunu/nunu/edit/feature/docs/:path',
      text: 'Edit this page on GitHub'
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2019-present Evan You'
    }
  }
})

function nav(): DefaultTheme.NavItem[] {
  return [
    {
      text: '首页',
      link: '/en/'
    },
    {
      text: '文档',
      link: '/en/getting-started',
      activeMatch: '/en/getting-started'
    },

    {
      text: pkg.version,
      items: [
        {
          text: '更新日志',
          link: 'https://github.com/go-nunu/nunu/blob/main/CHANGELOG.md'
        },
      ]
    }
  ]
}

function sidebarGuide(): DefaultTheme.SidebarItem[] {
  return [
    {
      text: '入门指引',
      collapsed: false,
      items: [
        {text: '引言', link: 'guide',},
        {text: '快速开始', link: 'getting-started',},
        {
          text: 'Nunu命令行工具', link: 'cli',
        },
      ]
    },
    {
      text: '基础概念',
      collapsed: false,
      items: [
        {text: '分层架构', link: 'architecture',},
        {
          text: 'Wire依赖注入', link: 'wire',
        },
      ]
    },
    {
      text: '核心组件',
      collapsed: false,
      items: [
        {
          text: 'Server', link: 'server',
        },
        {text: 'Handler', link: 'handler'},
        {
          text: 'Middleware', link: 'middleware',
        },
        {text: 'Service', link: 'service'},
        {
          text: 'repository',
          link: 'repository',
          items: [
            {text: '数据库', link: 'database'},
            {text: 'Redis', link: 'redis'},
          ]
        },
        {text: 'Model', link: 'model'},

        {text: 'Pkg', link: 'pkg'},
        {text: '日志', link: 'logger'},
        {text: '自动化文档', link: 'swagger'},
        {text: '单元测试', link: 'unit-test'},

      ]
    },
    {
      text: '构建部署',
      collapsed: false,
      items: [
        {text: 'PM2+Nginx', link: 'nginx'},
        {text: 'Docker', link: 'docker'},
        {text: 'Swarm', link: 'swarm'},
        {text: 'K8s', link: 'k8s'},
      ]
    },
    {
      text: '参考',
      collapsed: false,
      items: [
        {text: '贡献指南', link: 'pr'},
      ]
    },


    {text: '加入交流群', base: '/zh/reference/', link: 'site-config'}
  ]
}

export const search: DefaultTheme.AlgoliaSearchOptions['locales'] = {
  en: {
    placeholder: '搜索文档',
    translations: {
      button: {
        buttonText: '搜索文档',
        buttonAriaLabel: '搜索文档'
      },
      modal: {
        searchBox: {
          resetButtonTitle: '清除查询条件',
          resetButtonAriaLabel: '清除查询条件',
          cancelButtonText: '取消',
          cancelButtonAriaLabel: '取消'
        },
        startScreen: {
          recentSearchesTitle: '搜索历史',
          noRecentSearchesText: '没有搜索历史',
          saveRecentSearchButtonTitle: '保存至搜索历史',
          removeRecentSearchButtonTitle: '从搜索历史中移除',
          favoriteSearchesTitle: '收藏',
          removeFavoriteSearchButtonTitle: '从收藏中移除'
        },
        errorScreen: {
          titleText: '无法获取结果',
          helpText: '你可能需要检查你的网络连接'
        },
        footer: {
          selectText: '选择',
          navigateText: '切换',
          closeText: '关闭',
          searchByText: '搜索提供者'
        },
        noResultsScreen: {
          noResultsText: '无法找到相关结果',
          suggestedQueryText: '你可以尝试查询',
          reportMissingResultsText: '你认为该查询应该有结果？',
          reportMissingResultsLinkText: '点击反馈'
        }
      }
    }
  }
}
