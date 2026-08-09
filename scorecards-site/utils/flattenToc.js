export const flattenToc = (links = []) =>
  links.flatMap((link) => [
    { id: link.id, depth: link.depth, text: link.text },
    ...flattenToc(link.children),
  ])
