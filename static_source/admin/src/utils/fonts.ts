import api from "@/api/api";

export async function loadFonts(variableName: string) {
  const res = await api.v1.variableServiceGetVariableByName(variableName)
    .catch(() => {
    })
    .finally(() => {

    })
  const {data} = res;
  const customFont = new FontFace(variableName, `url(data:font/auto;base64,${data.value})`);
  customFont.load().then((font) => {
    document.fonts.add(font);
  });
}

