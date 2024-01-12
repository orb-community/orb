import { StatusCodes } from 'http-status-codes';

export function getStatusCodeText(statusCode: number): string {
    return StatusCodes[statusCode];
}
