import { css } from '@emotion/css';

import { GrafanaTheme2 } from '@grafana/data';

export const cardStyle = (theme: GrafanaTheme2, complete: boolean) => {
  const completeGradient = 'linear-gradient(to right, #c1d500 30%, #9a989b 60%, #ffffff 100%)';
  const darkThemeGradients = complete ? completeGradient : 'linear-gradient(to right, #c1d500 30%, #9a989b 60%, #ffffff 100%)';
  const lightThemeGradients = complete ? completeGradient : 'linear-gradient(to right, #c1d500 30%, #9a989b 60%, #ffffff 100%)';

  const borderGradient = theme.isDark ? darkThemeGradients : lightThemeGradients;

  return {
    backgroundColor: theme.colors.background.secondary,
    marginRight: theme.spacing(4),
    border: `1px solid ${theme.colors.border.weak}`,
    borderBottomLeftRadius: theme.shape.borderRadius(2),
    borderBottomRightRadius: theme.shape.borderRadius(2),
    position: 'relative',
    maxHeight: '230px',

    [theme.breakpoints.down('xxl')]: {
      marginRight: theme.spacing(2),
    },
    '&::before': {
      display: 'block',
      content: "' '",
      position: 'absolute',
      left: 0,
      right: 0,
      height: '2px',
      top: 0,
      backgroundImage: borderGradient,
    },
  } as const;
};

export const cardContent = css({
  padding: '16px',
});
