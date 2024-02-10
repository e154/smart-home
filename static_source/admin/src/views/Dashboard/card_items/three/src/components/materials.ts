import {
  BoxBufferGeometry,
  MathUtils,
  Mesh,
  MeshStandardMaterial,
  TextureLoader,
} from 'three';

const createMaterial = () => {
// create a texture loader.
  const textureLoader = new TextureLoader();

// load a texture
  const texture = textureLoader.load(
    '/assets/textures/uv-test-bw.png',
  );

  // create a "standard" material using
  // the texture we just loaded as a color map
  const material = new MeshStandardMaterial({
    map: texture,
  });

  return material;
}

export {createMaterial}
