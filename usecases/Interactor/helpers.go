package Interactor

func cutPostfix(text, postfix string) (shortCutText string) {
	lenPostfix := len(postfix)
	if len(text) > lenPostfix {
		startPostfix := len(text) - lenPostfix
		packageNamePostfix := text[startPostfix:]
		if packageNamePostfix == postfix {
			shortCutText = text[0:startPostfix]
		}
	}
	return
}
