const withMDX = require('@next/mdx')({
  extension: /\.mdx?$/,
  options: {
    remarkPlugins: [],
    rehypePlugins: [],
    providerImportSource: "@mdx-js/react", // Recommended for MDX v2+
  },
})
module.exports = withMDX({
  pageExtensions: ['js', 'jsx', 'md', 'mdx'],
  output: 'export',
  images: {
    unoptimized: true,
  },
})
