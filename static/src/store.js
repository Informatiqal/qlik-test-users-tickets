import { writable } from "svelte/store";

function createUser() {
  const { subscribe, set, update } = writable(undefined);

  return {
    subscribe,
    select: (value) => update(() => value),
    reset: () => set(undefined),
  };
}

export const selectedUser = createUser();

function createVP() {
  const { subscribe, set, update } = writable(undefined);

  return {
    subscribe,
    select: (value) => update(() => value),
    reset: () => set(undefined),
  };
}

export const selectedVP = createVP();

function createProxy() {
  const { subscribe, set, update } = writable(undefined);

  return {
    subscribe,
    select: (value) => update(() => value),
    reset: () => set(undefined),
  };
}

export const selectedProxy = createProxy();

function showHideAboutMethod() {
  const { subscribe, set, update } = writable(false);

  return {
    subscribe,
    set: (value) => update(() => value),
    reset: () => set(undefined),
  };
}

export const showHideAbout = showHideAboutMethod();
