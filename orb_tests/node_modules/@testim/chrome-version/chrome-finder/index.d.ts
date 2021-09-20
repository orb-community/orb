/**
 * find a executable chrome for all support system
 * @returns {string} executable chrome full path
 * @throws
 * if no executable chrome find, ERROR_NO_INSTALLATIONS_FOUND will be throw
 * if platform is not one if ['win32','darwin','linux'], ERROR_PLATFORM_NOT_SUPPORT will be throw
 */
declare function findChrome(): string;
export = findChrome;
