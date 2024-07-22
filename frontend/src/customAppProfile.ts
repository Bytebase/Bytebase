import { useActuatorV1Store, useAppFeature } from "./store";
import { useCustomTheme } from "./utils/customTheme";

export const overrideAppProfile = () => {
  const query = new URLSearchParams(window.location.search);
  const actuatorStore = useActuatorV1Store();
  const mode = query.get("mode");
  if (mode === "STANDALONE") {
    actuatorStore.appProfile.embedded = true;

    // mode=STANDALONE is not easy to read, but for legacy support we keep it as
    // some customers are using it.
    actuatorStore.overrideAppFeatures({
      "bb.feature.embedded-in-iframe": true,
      "bb.feature.hide-help": true,
      "bb.feature.hide-quick-start": true,
      "bb.feature.hide-release-remind": true,
      "bb.feature.disallow-share-worksheet": true,
      "bb.feature.disallow-navigate-to-console": true,
    });
  }

  const customTheme = query.get("customTheme");
  if (customTheme === "lixiang") {
    actuatorStore.overrideAppFeatures({
      "bb.feature.custom-query-datasource": true,
      "bb.feature.disallow-export-query-data": true,
      "bb.feature.custom-color-scheme": {
        "--color-accent": "#00665f",
        "--color-accent-hover": "#00554f",
        "--color-accent-disabled": "#b8c3c3",
      },
    });
    if (actuatorStore.appProfile.embedded) {
      actuatorStore.overrideAppFeatures({
        "bb.feature.hide-issue-review-actions": true,
      });
    }
  }

  useCustomTheme(useAppFeature("bb.feature.custom-color-scheme"));
};
