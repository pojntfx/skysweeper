import { IConfiguration } from "./models";

export class ConfigurationRestAPI {
  constructor(
    private apiURL: URL,
    private service: string,
    private accessJWT: string
  ) {}

  async getConfiguration(): Promise<IConfiguration> {
    const configurationURL = new URL(this.apiURL + "configuration");

    configurationURL.search = new URLSearchParams({
      service: this.service,
    }).toString();

    return (
      await fetch(configurationURL.toString(), {
        headers: {
          Authorization: "Bearer " + this.accessJWT,
          "Content-Type": "application/json",
        },
      })
    ).json();
  }
}
