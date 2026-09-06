export const useHeaderColour = () =>
  useState('headerColour', () => ({ bg: null, textColor: null }))

export const setHeaderColour = ({ bg, text }) => {
  const headerColour = useHeaderColour()
  headerColour.value.bg = bg
  headerColour.value.textColor = text
}
