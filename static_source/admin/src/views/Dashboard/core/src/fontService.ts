import api from "@/api/api";

class FontService {
  private loadedFonts: { [key: string]: FontFace } = {};

  removeFont = (variableName) => {
    if (!loadedFonts[variableName]) return
    document.fonts.delete(loadedFonts[variableName]);
  }

  loadFont = (variableName: string) => {
    api.v1.variableServiceGetVariableByName(variableName).then((res) => {
      const {data} = res;
      const customFont = new FontFace(variableName, `url(data:font/auto;base64,${data.value})`);
      customFont.load().then((font) => {
        this.loadedFonts[variableName] = font
        document.fonts.add(font);
      });
    })
  }
}

export const fontService = new FontService();
