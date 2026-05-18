// @ts-check

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'OTPLet homepage',
  tagline: 'Documentation',
  favicon: 'img/favicon.ico',

  // REQUIRED for routing
  url: 'https://otplet.cassano.fe.it',
  baseUrl: '/',

  // REQUIRED for GitHub Pages deploy
  organizationName: 'spetix',
  projectName: 'otplet',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.js',
        },
        blog: false, // disable blog
        theme: {
          customCss: './src/css/custom.css',
        },
      },
    ],
  ],
};

export default config;
