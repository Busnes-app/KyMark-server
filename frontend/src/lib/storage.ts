// Local device vault storage for zero-knowledge encryption keys
// Enables 1-click SSO unlock on trusted devices without compromising zero-knowledge server boundary.

const DB_NAME = 'kymark-device-vault';
const LEGACY_DB_NAME = 'kybookmarks-device-vault';
const STORE_NAME = 'keys';

function openDatabase(name = DB_NAME): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(name, 1);
    request.onupgradeneeded = () => {
      if (!request.result.objectStoreNames.contains(STORE_NAME)) {
        request.result.createObjectStore(STORE_NAME, { keyPath: 'username' });
      }
    };
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error ?? new Error('Unable to open local key store'));
  });
}

function deleteDatabase(name: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.deleteDatabase(name);
    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error ?? new Error('Unable to delete local key store'));
  });
}

export async function migrateLegacyDeviceVault(): Promise<void> {
  const legacy = await openDatabase(LEGACY_DB_NAME);
  const entries = await new Promise<Array<{ username: string; rawKey: string; updatedAt?: string }>>((resolve, reject) => {
    const request = legacy.transaction(STORE_NAME, 'readonly').objectStore(STORE_NAME).getAll();
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
  legacy.close();

  if (entries.length > 0) {
    const current = await openDatabase();
    await new Promise<void>((resolve, reject) => {
      const transaction = current.transaction(STORE_NAME, 'readwrite');
      for (const entry of entries) transaction.objectStore(STORE_NAME).put(entry);
      transaction.oncomplete = () => resolve();
      transaction.onerror = () => reject(transaction.error);
    });
    current.close();
  }
  await deleteDatabase(LEGACY_DB_NAME);
}

export async function storeDeviceVaultKey(username: string, rawKeyBase64: string): Promise<void> {
  const db = await openDatabase();
  await new Promise<void>((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readwrite');
    const store = transaction.objectStore(STORE_NAME);
    const request = store.put({
      username,
      rawKey: rawKeyBase64,
      updatedAt: new Date().toISOString(),
    });
    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error);
  });
  db.close();
}

export async function getDeviceVaultKey(username: string): Promise<string | undefined> {
  const db = await openDatabase();
  const result = await new Promise<{ username: string; rawKey: string } | undefined>((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readonly');
    const store = transaction.objectStore(STORE_NAME);
    const request = store.get(username);
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
  db.close();
  return result?.rawKey;
}

export async function clearDeviceVaultKey(username: string): Promise<void> {
  const db = await openDatabase();
  await new Promise<void>((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readwrite');
    const store = transaction.objectStore(STORE_NAME);
    const request = store.delete(username);
    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error);
  });
  db.close();
  await deleteDatabase(LEGACY_DB_NAME);
}

export async function clearAllDeviceVaultKeys(): Promise<void> {
  const db = await openDatabase();
  await new Promise<void>((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readwrite');
    const store = transaction.objectStore(STORE_NAME);
    const request = store.clear();
    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error);
  });
  db.close();
  await deleteDatabase(LEGACY_DB_NAME);
}
