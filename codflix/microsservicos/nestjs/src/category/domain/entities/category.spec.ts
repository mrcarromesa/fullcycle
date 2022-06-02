import { Category, CategoryProperties } from "./category";
import UniqueEntityId from "../../../@seedwork/domain/unique-entity-id.vo";

describe("Category Tests", () => {
  test("constructor of category", () => {
    jest.spyOn(Date, 'now').mockImplementation(() => new Date('2022-01-01 00:00:00').getTime());

    let category = new Category({
      name: "Movie",
      description: "Movie description",
      is_active: true,
    });
    expect(category.props).toStrictEqual({
      name: "Movie",
      description: "Movie description",
      is_active: true,
      created_at: new Date('2022-01-01 00:00:00'),
    });

    category = new Category({
      name: "Movie",
      description: "Movie description",
      is_active: false,
    });
    expect(category.props).toStrictEqual({
      name: "Movie",
      description: "Movie description",
      is_active: false,
      created_at: new Date('2022-01-01 00:00:00'),
    });
    
    category = new Category({
      name: "Movie",
      description: "Movie description",
    });
    expect(category.props).toMatchObject({
      name: "Movie",
      description: "Movie description",
    });
    
    category = new Category({
      name: "Movie",
      is_active: true,
    });
    expect(category.props).toMatchObject({
      name: "Movie",
      is_active: true,
    });
    
    category = new Category({
      name: "Movie",
      created_at: new Date('2022-01-01 00:00:00'),
    });
    expect(category.props).toMatchObject({
      name: "Movie",
      created_at: new Date('2022-01-01 00:00:00'),
    });

  });

  test('id field', () => {

    type CategoryData = { props: CategoryProperties, id?: UniqueEntityId };
    const data: CategoryData[] = [
      { props: { name: "Movie" }, },
      { props: { name: "Movie" }, id: null },
      { props: { name: "Movie" }, id: undefined },
      { props: { name: "Movie" }, id: new UniqueEntityId() },
    ]

    data.forEach(i => {
      const category = new Category(i.props, i.id);
      expect(category.id).not.toBeNull();
      expect(category.id).toBeInstanceOf(UniqueEntityId);
    });
  });

  test('getter of name field', () => {
    let category = new Category({
      name: "Movie",
      description: "Movie description",
      is_active: true,
    });
    expect(category.name).toBe("Movie");
  });
  
  test('getter and setter of description field', () => {
    let category = new Category({
      name: "Movie",
    });
    expect(category.description).toBeNull();
    
    category = new Category({
      name: "Movie",
      description: "Movie description",
    });
    expect(category.description).toBe("Movie description");
    
    category = new Category({
      name: "Movie",
    });
    // access private property
    category["description"] = "Movie description1";
    expect(category.description).toBe("Movie description1");
    
    category = new Category({
      name: "Movie",
      description: "Movie description",
    });
    category["description"] = undefined;
    expect(category.description).toBeNull();
    
    category = new Category({
      name: "Movie",
      description: "Movie description",
    });
    category["description"] = null;
    expect(category.description).toBeNull();
  });

  test('getter and setter of is_active field', () => {
    let category = new Category({
      name: "Movie",
    });
    expect(category.is_active).toBeTruthy();
    
    category = new Category({
      name: "Movie",
      is_active: false,
    });
    expect(category.is_active).toBeFalsy();
    
    category = new Category({
      name: "Movie",
    });
    // access private property
    category["is_active"] = false;
    expect(category.is_active).toBeFalsy();
    
    category = new Category({
      name: "Movie",
      is_active: true,
    });
    category["is_active"] = undefined;
    expect(category.is_active).toBeNull();
    
    category = new Category({
      name: "Movie",
      is_active: true,
    });
    category["is_active"] = null;
    expect(category.is_active).toBeNull();
  });

  test('getter and setter of created_at field', () => {
    let category = new Category({
      name: "Movie",
    });
    expect(category.created_at).toBeDefined();
    expect(category.created_at).toBeInstanceOf(Date);
    
    category = new Category({
      name: "Movie",
      created_at: new Date('2022-01-01 00:00:00'),
    });
    expect(category.created_at).toBeInstanceOf(Date);
  });
})