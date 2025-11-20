import type { Attribute, Schema } from '@strapi/strapi';

export interface RouteRouteStop extends Schema.Component {
  collectionName: 'components_route_route_stops';
  info: {
    description: '\u041E\u0441\u0442\u0430\u043D\u043E\u0432\u043A\u0430 \u043C\u0430\u0440\u0448\u0440\u0443\u0442\u0430 \u0441 \u043C\u0435\u0441\u0442\u043E\u043C \u0438 \u043F\u043E\u0440\u044F\u0434\u043A\u043E\u0432\u044B\u043C \u043D\u043E\u043C\u0435\u0440\u043E\u043C';
    displayName: '\u041E\u0441\u0442\u0430\u043D\u043E\u0432\u043A\u0430 \u043C\u0430\u0440\u0448\u0440\u0443\u0442\u0430';
    icon: 'map-marker-alt';
  };
  attributes: {
    order_num: Attribute.Integer &
      Attribute.Required &
      Attribute.SetMinMax<
        {
          min: 0;
        },
        number
      >;
    place: Attribute.Relation<
      'route.route-stop',
      'manyToOne',
      'api::place.place'
    > &
      Attribute.Required;
  };
}

declare module '@strapi/types' {
  export module Shared {
    export interface Components {
      'route.route-stop': RouteRouteStop;
    }
  }
}
